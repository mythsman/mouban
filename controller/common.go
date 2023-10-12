package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BizError(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, gin.H{
		"success": false,
		"msg":     msg,
	})
}
