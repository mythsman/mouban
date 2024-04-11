package agent

import (
	"bytes"
	"context"
	"io"
	"mouban/common"
	"mouban/model"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/sirupsen/logrus"
)

func TestDb(t *testing.T) {

	cfg := aws.NewConfig()
	cfg.BaseEndpoint = aws.String("https://s3.mythsman.com")
	cfg.Region = "us-east-1"
	cfg.Credentials = credentials.StaticCredentialsProvider{
		Value: aws.Credentials{
			AccessKeyID:     "RaE76H28Rk",
			SecretAccessKey: "_4UyUHn2bt.W6pNNHf@c",
		},
	}

	client := s3.NewFromConfig(*cfg)

	logrus.Infoln("test start")

	var storages []model.Storage
	id := 0
	step := 100
	for {
		wg := new(sync.WaitGroup)

		common.Db.Where("id > ? ", id).Order("id asc").Limit(step).Find(&storages)
		id += step

		for i := range storages {
			wg.Add(1)
			go func(s3Client *s3.Client, url string) {
				defer wg.Done()
				check(s3Client, url[strings.Index(url, "douban"):])
			}(client, storages[i].Target)
		}

		wg.Wait()
		if id%10000 == 0 {
			logrus.Infoln(id, "done")
		}
	}
}
func check(client *s3.Client, name string) {

	_, err := client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(""),
		Key:    aws.String(name),
	})

	if err != nil {
		logrus.Println(name, "check failed")
	}

}
func restore(client *s3.Client, url string, name string) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		logrus.Infoln("get failed for", url, name)
	}
	defer resp.Body.Close()

	contentType := resp.Header["Content-Type"][0]

	data, _ := io.ReadAll(resp.Body)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(""),
		Key:         aws.String(name),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})

	if err != nil {
		logrus.Warnln(name, "restore failed", err)
	}

}
