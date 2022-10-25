package main

import (
	"github.com/gin-gonic/gin"
	"mouban/controller"
)

func CollectRoute(router *gin.Engine) *gin.Engine {

	openGroup := router.Group("/open")
	{
		openGroup.GET("/user_book", controller.ListUserBook)
		openGroup.GET("/user_game", controller.ListUserGame)
		openGroup.GET("/user_movie", controller.ListUserMovie)
		openGroup.GET("/user_music", controller.ListUserMusic)

		openGroup.GET("/user", controller.GetUser)
		openGroup.GET("/book", controller.GetBook)
		openGroup.GET("/game", controller.GetGame)
		openGroup.GET("/movie", controller.GetMovie)
		openGroup.GET("/music", controller.GetMusic)
	}

	adminGroup := router.Group("/admin")
	{
		adminGroup.GET("/data_info", controller.GetDataInfo)
	}
	return router
}
