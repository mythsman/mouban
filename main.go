package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mouban/controller"
	"net/http"
)

func main() {

	router := gin.Default()

	router.Use(Recover)

	queryGroup := router.Group("/guest")
	{
		queryGroup.GET("/check_user", controller.CheckUser)
		queryGroup.GET("/user_book", controller.ListUserBook)
		queryGroup.GET("/user_game", controller.ListUserGame)
		queryGroup.GET("/user_movie", controller.ListUserMovie)
	}

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/overview", controller.GetOverview)
		adminGroup.GET("/crawl_user", controller.CrawlUser)
		adminGroup.GET("/crawl_book", controller.CrawlBook)
		adminGroup.GET("/crawl_movie", controller.CrawlMovie)
		adminGroup.GET("/crawl_game", controller.CrawlGame)
	}

	panic(router.Run(":" + viper.GetString("server.port")))
}

func Recover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"success": false,
				"msg":     r,
			})
		}
	}()
	ctx.Next()
}
