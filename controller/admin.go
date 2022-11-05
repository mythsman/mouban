package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mouban/consts"
	"mouban/dao"
	"strconv"
)

func GetOverview(ctx *gin.Context) {
	if !checkPermission(ctx) {
		return
	}
}

func CrawlUser(ctx *gin.Context) {
	if !checkPermission(ctx) {
		return
	}
	id := ctx.Query("id")
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		BizError(ctx, "参数错误")
		return
	}
	dispatch(idInt, consts.TypeUser)
}

func CrawlItem(ctx *gin.Context, t uint8) {
	if !checkPermission(ctx) {
		return
	}
	id := ctx.Query("id")
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		BizError(ctx, "参数错误")
		return
	}
	dispatch(idInt, t)
}

func checkPermission(ctx *gin.Context) bool {
	token := ctx.Query("token")
	if token != viper.GetString("admin.token") {
		BizError(ctx, "参数错误")
		return false
	}
	return true
}

func dispatch(doubanId uint64, t uint8) bool {
	schedule := dao.GetSchedule(doubanId, t)
	triggered := false
	switch schedule.Status {
	case consts.ScheduleStatusCrawled:
		dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
		triggered = true
		break
	case consts.ScheduleStatusCrawling:
		break
	case consts.ScheduleStatusToCrawl:
		break
	default:
		break
	}
	return triggered
}
