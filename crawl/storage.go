package crawl

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/gabriel-vasile/mimetype"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Storage source url -> stored url
func Storage(url string) string {
	file := download(url, "https://www.douban.com/")
	mtype, extension := mime(file.Name())

	md5Result := md5sum(file.Name())

	result := upload(file.Name(), md5Result+extension, mtype)

	e := os.Remove(file.Name())

	if e != nil {
		log.Fatal(e)
	}
	return result
}

func download(url string, referer string) *os.File {
	// 创建一个文件用于保存
	out, err := os.CreateTemp("/tmp", "mouban-")
	if err != nil {
		panic(err)
	}
	defer out.Close()

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Referer", referer)

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
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
	endpoint := viper.GetString("minio.endpoint")
	accessKeyID := viper.GetString("minio.id")
	secretAccessKey := viper.GetString("minio.key")

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatalln(err)
	}
	bucketName := "douban"

	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	options := minio.PutObjectOptions{
		ContentType: mimeType,
	}
	_, err = minioClient.FPutObject(context.Background(), bucketName, name, file, options)
	if err != nil {
		log.Fatalln(err)
	}
	return "https://" + endpoint + "/" + bucketName + "/" + name
}
