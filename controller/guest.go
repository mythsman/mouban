package controller

import (
	"github.com/gin-gonic/gin"
	"mouban/consts"
	"mouban/dao"
	"net/http"
	"strconv"
)

func CheckUser(ctx *gin.Context) {
	id := ctx.Query("id")
	doubanUid, err := strconv.ParseUint(id, 10, 64)
	if err != nil || id == "0" {
		BizError(ctx, "用户ID输入错误")
		return
	}
	logAccess(ctx, doubanUid)

	schedule := dao.GetSchedule(doubanUid, consts.TypeUser)

	if schedule == nil {
		dao.CreateSchedule(doubanUid, consts.TypeUser, consts.ScheduleStatusToCrawl, consts.ScheduleResultUnready)
		BizError(ctx, "未录入当前用户，已发起录入，请等待后台数据更新")
		return
	}

	if schedule.Result == consts.ScheduleResultUnready {
		BizError(ctx, "当前用户录入中")
		return
	}

	if schedule.Result == consts.ScheduleResultInvalid {
		BizError(ctx, "当前用户不存在")
		return
	}

	user := dao.GetUser(doubanUid)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  user.Show(),
	})

}

func ListUserMovie(ctx *gin.Context) {
	id := ctx.Query("id")
	doubanUid, err := strconv.ParseUint(id, 10, 64)
	if err != nil || id == "0" {
		BizError(ctx, "用户ID输入错误")
		return
	}
	logAccess(ctx, doubanUid)

	action := ctx.Query("action")
	if action == "" {
		BizError(ctx, "参数错误")
		return
	}

	offset := 0
	if ctx.Query("offset") != "" {
		offset, _ = strconv.Atoi(ctx.Query("offset"))
	}

	schedule := dao.GetSchedule(doubanUid, consts.TypeUser)

	if schedule == nil {
		BizError(ctx, "当前用户未录入")
		return
	}

	if schedule.Result == consts.ScheduleResultUnready {
		BizError(ctx, "当前用户录入中")
		return
	}

	if schedule.Result == consts.ScheduleResultInvalid {
		BizError(ctx, "当前用户不存在")
		return
	}

	user := dao.GetUser(doubanUid)

	comments := dao.SearchComment(doubanUid, consts.TypeMovie, parseAction(action), offset, 20)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result": gin.H{
			"user":     user,
			"comments": comments,
		},
	})
}

func ListUserBook(ctx *gin.Context) {
	id := ctx.Query("id")
	doubanUid, err := strconv.ParseUint(id, 10, 64)
	if err != nil || id == "0" {
		BizError(ctx, "用户ID输入错误")
		return
	}
	logAccess(ctx, doubanUid)

	action := ctx.Query("action")
	if action == "" {
		BizError(ctx, "参数错误")
		return
	}
	schedule := dao.GetSchedule(doubanUid, consts.TypeUser)

	if schedule == nil {
		BizError(ctx, "当前用户未录入")
		return
	}

	if schedule.Result == consts.ScheduleResultUnready {
		BizError(ctx, "当前用户录入中")
		return
	}

	if schedule.Result == consts.ScheduleResultInvalid {
		BizError(ctx, "当前用户不存在")
		return
	}

	//user := dao.GetUser(doubanUid)

}

func ListUserGame(ctx *gin.Context) {
	id := ctx.Query("id")
	doubanUid, err := strconv.ParseUint(id, 10, 64)
	if err != nil || id == "0" {
		BizError(ctx, "用户ID输入错误")
		return
	}
	logAccess(ctx, doubanUid)

	action := ctx.Query("action")
	if action == "" {
		BizError(ctx, "参数错误")
		return
	}

	schedule := dao.GetSchedule(doubanUid, consts.TypeUser)

	if schedule == nil {
		BizError(ctx, "当前用户未录入")
		return
	}

	if schedule.Result == consts.ScheduleResultUnready {
		BizError(ctx, "当前用户录入中")
		return
	}

	if schedule.Result == consts.ScheduleResultInvalid {
		BizError(ctx, "当前用户不存在")
		return
	}

	//user := dao.GetUser(doubanUid)
}

func logAccess(ctx *gin.Context, doubanUid uint64) {
	ua := ctx.GetHeader("User-Agent")
	referer := ctx.GetHeader("Referer")
	ip := ctx.RemoteIP()

	dao.AddAccess(doubanUid, ctx.FullPath(), ip, ua, referer)
}

func parseAction(action string) uint8 {
	switch action {
	case consts.ActionWish.Name:
		return consts.ActionWish.Code
	case consts.ActionCollect.Name:
		return consts.ActionCollect.Code
	case consts.ActionDo.Name:
		return consts.ActionDo.Code
	}
	return consts.ActionCollect.Code
}
