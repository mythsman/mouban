package controller

import (
	"github.com/gin-gonic/gin"
	"mouban/consts"
	"mouban/logic"
	"strconv"
)

func GetOverview(ctx *gin.Context) {
}

func CrawlUser(ctx *gin.Context) {
	id := ctx.Query("id")
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		BizError(ctx, "参数错误")
		return
	}
	logic.Dispatch(idInt, consts.TypeUser)
}

func CrawlGame(ctx *gin.Context) {
	id := ctx.Query("id")
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		BizError(ctx, "参数错误")
		return
	}
	logic.Dispatch(idInt, consts.TypeGame)
}

func CrawlBook(ctx *gin.Context) {
	id := ctx.Query("id")
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		BizError(ctx, "参数错误")
		return
	}
	logic.Dispatch(idInt, consts.TypeBook)
}

func CrawlMovie(ctx *gin.Context) {
	id := ctx.Query("id")
	idInt, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		BizError(ctx, "参数错误")
		return
	}
	logic.Dispatch(idInt, consts.TypeMovie)
}
