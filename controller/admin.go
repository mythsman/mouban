package controller

import (
	"bufio"
	"context"
	"mouban/consts"
	"mouban/dao"
	"mouban/util"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

func RefreshUser(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id := util.ParseNumber(idStr)
	if id == 0 {
		BizError(ctx, "参数错误")
		return
	}
	schedule := dao.GetSchedule(id, consts.TypeUser.Code)
	if schedule == nil {
		BizError(ctx, "条目未收录，无法更新")
		return
	}

	user := dao.GetUser(id)
	if user == nil {
		BizError(ctx, "用户未收录，无法更新")
		return
	}

	dao.RefreshUser(user)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})

	dao.CasScheduleStatus(user.DoubanUid, consts.TypeUser.Code, consts.ScheduleToCrawl.Code, *schedule.Status)

}
func RefreshItem(ctx *gin.Context) {
	typeStr := ctx.Query("type")
	idStr := ctx.Query("id")

	t := uint8(util.ParseNumber(typeStr))
	id := util.ParseNumber(idStr)
	if t == 0 || id == 0 {
		BizError(ctx, "参数错误")
		return
	}

	schedule := dao.GetSchedule(id, t)
	if schedule == nil {
		BizError(ctx, "条目未收录，无法更新")
		return
	}

	if *schedule.Status == consts.ScheduleCrawling.Code || *schedule.Status == consts.ScheduleToCrawl.Code {
		BizError(ctx, "当前条目正在更新中")
		return
	}

	dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleToCrawl.Code, *schedule.Status)
	logrus.Infoln("refresh item for", t, id)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

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
		logrus.Errorln(err)
		return
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	sem := semaphore.NewWeighted(20)
	for scanner.Scan() {
		err := sem.Acquire(context.Background(), 1)
		if err != nil {
			logrus.Infoln("acquire semaphore failed", err)
			return
		}
		line := scanner.Text()
		go func() {
			defer func() {
				sem.Release(1)
			}()
			processLine(line)
		}()
	}

	if err := scanner.Err(); err != nil {
		logrus.Errorln(err)
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
