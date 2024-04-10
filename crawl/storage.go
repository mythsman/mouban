package crawl

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"mouban/dao"
	"mouban/model"
	"mouban/util"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_credentials "github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gabriel-vasile/mimetype"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

var minioClient *minio.Client
var retryClient *retryablehttp.Client
var s3Client *s3.Client

var endpoint string
var accessKeyID string
var secretAccessKey string
var bucketName string

// Storage source url -> stored url
func Storage(url string) string {

	if strings.Contains(url, viper.GetString("minio.endpoint")) {
		logrus.Infoln("storage ignore :", url)
		return url
	}

	if !strings.HasPrefix(url, "http") {
		logrus.Infoln("storage bad :", url)
		return ""
	}

	storageHit := dao.GetStorage(url)
	if storageHit != nil {
		logrus.Infoln("storage hit")
		return storageHit.Target
	}

	var file *os.File
	for i := 0; i < 5; i++ {
		file = download(url, "https://www.douban.com/")
		if file != nil {
			break
		}
	}
	if file == nil {
		panic("download file finally failed for : " + url)
	}

	mtype, extension := mime(file.Name())

	md5Result := md5sum(file.Name())

	result := ""
	existingStorage := dao.GetStorageByMd5(md5Result)
	if existingStorage != nil {
		result = existingStorage.Target
		logrus.Infoln("storage already uploaded for", md5Result)
	} else {
		result = upload(file.Name(), md5Result+extension, mtype)
	}

	_ = os.Remove(file.Name())

	storage := &model.Storage{
		Source: url,
		Target: result,
		Md5:    md5Result,
	}
	dao.UpsertStorage(storage)
	logrus.Infoln("storage add :", url, "->", result)

	if strings.HasSuffix(storage.Target, ".txt") || strings.HasSuffix(storage.Target, ".html") {
		logrus.Warnln("storage maybe invalid :", url, "->", result)
	}

	return result
}

func download(url string, referer string) (o *os.File) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("download panic", url, r, "=>", util.GetCurrentGoroutineStack())
			o = nil
		}
	}()
	// 创建一个文件用于保存
	out, err := os.CreateTemp("/tmp", "mouban-")
	if err != nil {
		logrus.Errorln("create tmp file failed")
		panic(err)
	}
	defer out.Close()

	req, _ := retryablehttp.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Referer", referer)

	resp, _ := retryClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		logrus.Errorln("write file", url, "failed")
		panic(err)
	}
	return out
}

func mime(path string) (string, string) {
	mtype, _ := mimetype.DetectFile(path)
	return mtype.String(), mtype.Extension()
}

func md5sum(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	hash := md5.New()
	_, _ = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil))
}

func upload(file string, name string, mimeType string) string {
	options := minio.PutObjectOptions{
		ContentType: mimeType,
	}
	_, err := minioClient.FPutObject(context.Background(), bucketName, name, file, options)
	if err != nil {
		logrus.Errorln("minio put failed,", err)
	}
	url := "https://" + endpoint + "/" + bucketName + "/" + name

	restore(url, bucketName+"/"+name)

	return url

}

func init() {
	retryClient = initHttpClient()
	endpoint = viper.GetString("minio.endpoint")
	accessKeyID = viper.GetString("minio.id")
	secretAccessKey = viper.GetString("minio.key")
	bucketName = viper.GetString("minio.bucket")

	// Initialize minio client object.
	err := errors.New("")
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})

	if err != nil {
		panic(err)
	}

	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			logrus.Println("We already own bucket", bucketName)
		}
	} else {
		logrus.Println("Successfully created bucket", bucketName)
	}

	initS3Client()
}

func initHttpClient() *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.RetryMax = viper.GetInt("http.retry_max")
	client.Logger = nil
	client.RetryWaitMin = time.Duration(1) * time.Second
	client.RetryWaitMax = time.Duration(60) * time.Second
	client.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		return retryablehttp.DefaultRetryPolicy(ctx, resp, err)
	}

	client.HTTPClient = &http.Client{
		Timeout: time.Duration(viper.GetInt("http.timeout")) * time.Millisecond,
	}
	return client
}

func initS3Client() {
	cfg := aws.NewConfig()
	cfg.BaseEndpoint = aws.String(viper.GetString("s3.endpoint"))
	cfg.Region = viper.GetString("s3.region")
	cfg.Credentials = aws_credentials.StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID:     viper.GetString("s3.access_key"),
			SecretAccessKey: viper.GetString("s3.secret_key"),
		},
	}

	s3Client = s3.NewFromConfig(*cfg)
}

func restore(url string, name string) {

	resp, err := http.Get(url)
	if err != nil {
		logrus.Infoln("get failed for", url, name)
	}
	defer resp.Body.Close()

	contentType := resp.Header["Content-Type"][0]

	data, _ := io.ReadAll(resp.Body)

	output, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(""),
		Key:         aws.String(name),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		logrus.Warnln(name, "restore failed", err)
	}
	logrus.Println(name, "restore done", contentType, *output.ETag)

}
