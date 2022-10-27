package dao

import (
	"github.com/gin-gonic/gin"
	"mouban/common"
	"mouban/model/do"
	"strconv"
)

func GetUser(ctx *gin.Context) {
	result := do.User{}

	doubanUidStr := ctx.Query("douban_uid")
	domain := ctx.Query("domain")
	if doubanUidStr != "" {
		doubanUid, _ := strconv.Atoi(doubanUidStr)

		common.Db.Where(&do.User{DoubanUid: uint64(doubanUid)}).Take(&result)
		ctx.JSON(200, result)
		return
	}

	if domain != "" {
		common.Db.Where(&do.User{Domain: domain}).Take(&result)
		ctx.JSON(200, result)
		return
	}
	ctx.JSON(200, result)
}
