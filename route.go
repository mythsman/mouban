package main

import (
	"github.com/gin-gonic/gin"
	"mouban/controller"
)

func CollectRoute(router *gin.Engine) *gin.Engine {

	queryGroup := router.Group("/guest")
	{
		queryGroup.GET("/user_book", controller.ListUserBook)
		queryGroup.GET("/user_game", controller.ListUserGame)
		queryGroup.GET("/user_movie", controller.ListUserMovie)
		queryGroup.GET("/user_music", controller.ListUserMusic)

		queryGroup.GET("/user", controller.GetUser)
		queryGroup.GET("/book", controller.GetBook)
		queryGroup.GET("/game", controller.GetGame)
		queryGroup.GET("/movie", controller.GetMovie)
		queryGroup.GET("/music", controller.GetMusic)
	}

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/info", controller.GetOverview)
	}

	volunteerGroup := router.Group("/volunteer")
	{
		volunteerGroup.GET("/pull_task", controller.Test)
	}
	return router
}
