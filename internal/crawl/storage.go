package crawl

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mouban/internal/dao"
	"mouban/internal/model"
	"mouban/internal/util"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gabriel-vasile/mimetype"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

var s3Client *s3.Client

var storageInitOnce sync.Once
var storageInitErr error

func initStorageFromConfig() error {
	endpoint := strings.TrimSpace(viper.GetString("s3.endpoint"))
	bucket := strings.TrimSpace(viper.GetString("s3.bucket"))
	region := strings.TrimSpace(viper.GetString("s3.region"))
	accessKey := strings.TrimSpace(viper.GetString("s3.access_key"))
	secretKey := strings.TrimSpace(viper.GetString("s3.secret_key"))

	if endpoint == "" || bucket == "" || region == "" || accessKey == "" || secretKey == "" {
		return fmt.Errorf("s3 config missing: require s3.endpoint/s3.bucket/s3.region/s3.access_key/s3.secret_key")
	}
	if _, err := url.ParseRequestURI(endpoint); err != nil {
		return fmt.Errorf("invalid s3.endpoint: %w", err)
	}

	s3Client = initS3Client()
	return nil
}

func ensureStorageInitialized() error {
	storageInitOnce.Do(func() {
		storageInitErr = initStorageFromConfig()
	})
	return storageInitErr
}

// Storage source url -> stored url
func Storage(url string) string {
	if !isStorageEnabled() {
		return url
	}

	if err := ensureStorageInitialized(); err != nil {
		panic(err)
	}

	if strings.Contains(url, viper.GetString("s3.endpoint")) {
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

const browserUA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"

func download(url string, referer string) (o *os.File) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("download panic", url, r, "=>", util.GetCurrentGoroutineStack())
			o = nil
		}
	}()

	out, err := os.CreateTemp("/tmp", "mouban-")
	if err != nil {
		logrus.Errorln("create tmp file failed")
		panic(err)
	}
	defer out.Close()

	statusCode, err := downloadWithCurl(url, referer, out.Name())
	if err == nil && isImageFile(out.Name()) {
		return out
	}

	if statusCode == 404 || statusCode == 418 {
		if truncateErr := out.Truncate(0); truncateErr != nil {
			panic(fmt.Errorf("truncate fallback file failed: %w", truncateErr))
		}
		logrus.Warnln("download fallback to empty file for status", statusCode, ":", url)
		return out
	}

	if err != nil {
		logrus.Warnln("download by curl failed:", url, err, "status", statusCode)
	} else {
		logrus.Warnln("download by curl got non-image:", url, "status", statusCode)
	}

	panic("download got invalid image for: " + url)
}

func downloadWithCurl(url string, referer string, output string) (int, error) {
	args := []string{"-L", "-sS", "-A", browserUA, "-e", referer, "-o", output, "-w", "%{http_code}", url}
	cmd := exec.Command("curl", args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	statusText := strings.TrimSpace(stdout.String())
	statusCode, _ := strconv.Atoi(statusText)
	if err != nil {
		return statusCode, fmt.Errorf("curl error: %w, stderr: %s", err, strings.TrimSpace(stderr.String()))
	}
	return statusCode, nil
}

func isImageFile(path string) bool {
	mtype, _ := mime(path)
	return strings.HasPrefix(mtype, "image/")
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
	f, _ := os.Open(file)
	defer f.Close()

	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(viper.GetString("s3.bucket")),
		Key:         aws.String(name),
		Body:        f,
		ContentType: aws.String(mimeType),
	})

	if err != nil {
		logrus.Warnln(name, "restore failed", err)
	}

	url := viper.GetString("s3.endpoint") + "/" + viper.GetString("s3.bucket") + "/" + name

	return url
}

func initS3Client() *s3.Client {
	cfg := aws.NewConfig()
	cfg.BaseEndpoint = aws.String(viper.GetString("s3.endpoint"))
	cfg.Region = viper.GetString("s3.region")
	cfg.Credentials = credentials.StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID:     viper.GetString("s3.access_key"),
			SecretAccessKey: viper.GetString("s3.secret_key"),
		},
	}

	return s3.NewFromConfig(*cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
}
