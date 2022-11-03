package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BizError(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": false,
		"msg":     msg,
	})
}
