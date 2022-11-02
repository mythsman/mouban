package controller

import (
	"github.com/gin-gonic/gin"
	"mouban/dao"
	"mouban/logic"
	"net/http"
	"strconv"
)

func CheckUser(ctx *gin.Context) {
	doubanUid := LogAccess(ctx)

	user := dao.GetUser(doubanUid)
	if user == nil {
		logic.DispatchUser(doubanUid)
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "未录入当前用户，请等待后台数据更新（约十分钟）",
		})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  user.Show(),
	})

}

func ListUserMovie(ctx *gin.Context) {
	doubanUid := LogAccess(ctx)

	user := dao.GetUser(doubanUid)
	if user == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "未知用户",
		})
		return
	}
	action := ctx.Query("action")
	if action == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "参数错误",
		})
		return
	}

}

func ListUserBook(ctx *gin.Context) {
	doubanUid := LogAccess(ctx)

	user := dao.GetUser(doubanUid)
	if user == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "未知用户",
		})
		return
	}
	action := ctx.Query("action")
	if action == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "参数错误",
		})
		return
	}

}

func ListUserGame(ctx *gin.Context) {
	doubanUid := LogAccess(ctx)

	user := dao.GetUser(doubanUid)
	if user == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "未知用户",
		})
		return
	}
	action := ctx.Query("action")
	if action == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "参数错误",
		})
		return
	}
}

func LogAccess(ctx *gin.Context) uint64 {
	id := ctx.Query("id")
	doubanUid, err := strconv.ParseUint(id, 10, 64)
	if err != nil || id == "0" {
		return 0
	}

	ua := ctx.GetHeader("User-Agent")
	referer := ctx.GetHeader("Referer")
	ip := ctx.RemoteIP()

	dao.AddAccess(doubanUid, ctx.FullPath(), ip, ua, referer)

	return doubanUid
}
