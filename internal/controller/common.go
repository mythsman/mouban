package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BizError(ctx *gin.Context, code int, msg string) {
	ctx.JSON(code, gin.H{
		"success": false,
		"msg":     msg,
	})
}

func BadRequest(ctx *gin.Context, msg string) {
	BizError(ctx, http.StatusBadRequest, msg)
}

func NotFound(ctx *gin.Context, msg string) {
	BizError(ctx, http.StatusNotFound, msg)
}

func Conflict(ctx *gin.Context, msg string) {
	BizError(ctx, http.StatusConflict, msg)
}

func Accepted(ctx *gin.Context, msg string) {
	BizError(ctx, http.StatusAccepted, msg)
}
