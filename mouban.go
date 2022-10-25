package main

import (
	"github.com/gin-gonic/gin"
	"mouban/common"
)

func init() {
	common.InitDB()
}
func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := r.Run()
	if err != nil {
		return
	}
}
