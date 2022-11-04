package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mouban/consts"
	"mouban/logic"
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
	logic.Dispatch(idInt, consts.TypeUser)
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
	logic.Dispatch(idInt, t)
}

func checkPermission(ctx *gin.Context) bool {
	token := ctx.Query("token")
	if token != viper.GetString("admin.token") {
		BizError(ctx, "参数错误")
		return false
	}
	return true
}
