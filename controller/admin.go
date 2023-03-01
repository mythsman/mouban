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

	logrus.Infoln("start loading ", path)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})

	go loadFile(path)
}

func loadFile(path string) {
	f, err := os.Open(path)

	if err != nil {
		logrus.Fatalln(err)
		return
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	sem := semaphore.NewWeighted(100)
	for scanner.Scan() {
		err := sem.Acquire(context.Background(), 1)
		if err != nil {
			logrus.Infoln("acquire semaphore failed", err)
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
		logrus.Fatalln(err)
	}
}

func processLine(line string) {
	doubanId, t := util.ParseItem(line)
	if doubanId == 0 {
		return
	}
	schedule := dao.GetSchedule(doubanId, t.Code)
	if schedule == nil {
		added := dao.CreateScheduleNx(doubanId, t.Code, consts.ScheduleCanCrawl.Code, consts.ScheduleUnready.Code)
		if added {
			logrus.Infoln("new", t.Name, "added :", doubanId)
		} else {
			logrus.Infoln("new", t.Name, "duplicated :", doubanId)
		}
	} else {
		logrus.Infoln("old", t.Name, "ignored :", doubanId)
	}
}
