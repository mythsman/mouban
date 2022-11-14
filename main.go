package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"mouban/consts"
	"mouban/controller"
	_ "mouban/routine"
	"mouban/util"
	"net/http"
)

func main() {

	router := gin.Default()

	router.Use(Recover)

	queryGroup := router.Group("/guest")
	{
		queryGroup.GET("/check_user", controller.CheckUser)
		queryGroup.GET("/user_book", func(ctx *gin.Context) {
			controller.ListUserItem(ctx, consts.TypeBook)
		})
		queryGroup.GET("/user_game", func(ctx *gin.Context) {
			controller.ListUserItem(ctx, consts.TypeGame)
		})
		queryGroup.GET("/user_movie", func(ctx *gin.Context) {
			controller.ListUserItem(ctx, consts.TypeMovie)
		})
	}

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/overview", controller.GetOverview)
		adminGroup.GET("/crawl_user", controller.CrawlUser)
		adminGroup.GET("/crawl_book", func(ctx *gin.Context) {
			controller.CrawlItem(ctx, consts.TypeBook)
		})
		adminGroup.GET("/crawl_movie", func(ctx *gin.Context) {
			controller.CrawlItem(ctx, consts.TypeMovie)
		})
		adminGroup.GET("/crawl_game", func(ctx *gin.Context) {
			controller.CrawlItem(ctx, consts.TypeGame)
		})
	}

	panic(router.Run(":" + viper.GetString("server.port")))
}

func Recover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"msg":     "服务内部错误，请联系开发者处理",
			})
			log.Println(r, " => ", util.GetCurrentGoroutineStack())
		}
	}()
	ctx.Next()
}
