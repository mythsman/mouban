package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mouban/controller"
)

func main() {

	router := gin.Default()

	router = collectRoute(router)

	panic(router.Run(":" + viper.GetString("server.port")))
}

func collectRoute(router *gin.Engine) *gin.Engine {

	queryGroup := router.Group("/guest")
	{
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

	volunteerGroup := router.Group("/volunteer")
	{
		volunteerGroup.POST("/report_user", controller.ReportUser)
	}
	return router
}
