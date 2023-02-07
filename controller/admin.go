package controller

import (
	"bufio"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"mouban/consts"
	"mouban/dao"
	"mouban/util"
	"net/http"
	"os"
)

func LoadData(ctx *gin.Context) {
	path := ctx.Query("path")

	logrus.Info("start loading ", path)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})

	go loadFile(path)
}

func loadFile(path string) {
	f, err := os.Open(path)

	if err != nil {
		logrus.Fatal(err)
		return
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	sem := semaphore.NewWeighted(100)
	for scanner.Scan() {
		err := sem.Acquire(context.Background(), 1)
		if err != nil {
			logrus.Info("acquire semaphore failed", err)
			return
		}
		go func() {
			defer func() {
				sem.Release(1)
			}()
			processLine(scanner.Text())
		}()
	}

	if err := scanner.Err(); err != nil {
		logrus.Fatal(err)
	}
}

func processLine(line string) {
	doubanId, t := util.ParseItem(line)
	if doubanId == 0 {
		return
	}
	schedule := dao.GetSchedule(doubanId, t.Code)
	if schedule == nil {
		added := dao.CreateScheduleNx(doubanId, t.Code, consts.ScheduleStatusCanCrawl, consts.ScheduleResultUnready)
		if added {
			logrus.Info("new", t.Name, "added :", doubanId)
		} else {
			logrus.Info("new", t.Name, "duplicated :", doubanId)
		}
	} else {
		logrus.Info("old", t.Name, "ignored :", doubanId)
	}
}
