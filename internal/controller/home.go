package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HomePage(ctx *gin.Context) {
	logAccess(ctx, 0)
	ctx.HTML(http.StatusOK, "home.tmpl", gin.H{})
}

func ExplorePage(ctx *gin.Context) {
	logAccess(ctx, 0)
	ctx.HTML(http.StatusOK, "explore.tmpl", gin.H{})
}
