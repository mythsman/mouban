package controller

import (
	"github.com/gin-gonic/gin"
	"mouban/common"
	"mouban/model"
	"strconv"
)

func ListUserMusic(ctx *gin.Context) {

}
func ListUserMovie(ctx *gin.Context) {

}
func ListUserBook(ctx *gin.Context) {

}
func ListUserGame(ctx *gin.Context) {

}

func GetUser(ctx *gin.Context) {
	result := model.User{}

	doubanUidStr := ctx.Query("douban_uid")
	uniqueId := ctx.Query("unique_id")
	if doubanUidStr != "" {
		doubanUid, _ := strconv.Atoi(doubanUidStr)

		common.Db.Where(&model.User{DoubanUid: uint64(doubanUid)}).Take(&result)
		ctx.JSON(200, result)
		return
	}

	if uniqueId != "" {
		common.Db.Where(&model.User{UniqueId: uniqueId}).Take(&result)
		ctx.JSON(200, result)
		return
	}
	ctx.JSON(200, result)
}

func GetBook(ctx *gin.Context) {

}

func GetMovie(ctx *gin.Context) {

}

func GetMusic(ctx *gin.Context) {

}
func GetGame(ctx *gin.Context) {

}
