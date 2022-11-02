package controller

import (
	"github.com/gin-gonic/gin"
	"mouban/consts"
	"mouban/dao"
	"mouban/logic"
	"net/http"
	"strconv"
)

func CheckUser(ctx *gin.Context) {
	doubanUid := logAccess(ctx)

	schedule := dao.GetSchedule(doubanUid, consts.TypeUser)

	if schedule == nil {
		logic.Dispatch(doubanUid, consts.TypeUser)
		panic("未录入当前用户，已发起录入，请等待后台数据更新")
	}

	if schedule.Result == consts.ScheduleResultUnready {
		panic("当前用户录入中")
	}

	if schedule.Result == consts.ScheduleResultInvalid {
		panic("当前用户不存在")
	}

	user := dao.GetUser(doubanUid)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  user.Show(),
	})

}

func ListUserMovie(ctx *gin.Context) {
	doubanUid := logAccess(ctx)

	action := ctx.Query("action")
	if action == "" {
		panic("参数错误")
	}
	schedule := dao.GetSchedule(doubanUid, consts.TypeUser)

	if schedule == nil {
		panic("当前用户未录入")
	}

	if schedule.Result == consts.ScheduleResultUnready {
		panic("当前用户录入中")
	}

	if schedule.Result == consts.ScheduleResultInvalid {
		panic("当前用户不存在")
	}

	//user := dao.GetUser(doubanUid)

}

func ListUserBook(ctx *gin.Context) {
	doubanUid := logAccess(ctx)

	action := ctx.Query("action")
	if action == "" {
		panic("参数错误")
	}
	schedule := dao.GetSchedule(doubanUid, consts.TypeUser)

	if schedule == nil {
		panic("当前用户未录入")
	}

	if schedule.Result == consts.ScheduleResultUnready {
		panic("当前用户录入中")
	}

	if schedule.Result == consts.ScheduleResultInvalid {
		panic("当前用户不存在")
	}

	//user := dao.GetUser(doubanUid)

}

func ListUserGame(ctx *gin.Context) {
	doubanUid := logAccess(ctx)

	action := ctx.Query("action")
	if action == "" {
		panic("参数错误")
	}

	schedule := dao.GetSchedule(doubanUid, consts.TypeUser)

	if schedule == nil {
		panic("当前用户未录入")
	}

	if schedule.Result == consts.ScheduleResultUnready {
		panic("当前用户录入中")
	}

	if schedule.Result == consts.ScheduleResultInvalid {
		panic("当前用户不存在")
	}

	//user := dao.GetUser(doubanUid)
}

func logAccess(ctx *gin.Context) uint64 {
	id := ctx.Query("id")
	doubanUid, err := strconv.ParseUint(id, 10, 64)
	if err != nil || id == "0" {
		panic("用户ID输入错误")
	}

	ua := ctx.GetHeader("User-Agent")
	referer := ctx.GetHeader("Referer")
	ip := ctx.RemoteIP()

	dao.AddAccess(doubanUid, ctx.FullPath(), ip, ua, referer)

	return doubanUid
}
