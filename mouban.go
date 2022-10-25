package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mouban/common"
)

func init() {
	common.InitDB()
}

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	infoGroup := router.Group("/info")
	{
		infoGroup.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "test",
			})
		})
	}

	actionGroup := router.Group("/action")
	{
		actionGroup.GET("/get", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "get",
			})
		})
	}

	err := router.Run(":" + viper.GetString("server.port"))
	if err != nil {
		return
	}
}
