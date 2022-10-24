package main

import (
	"github.com/gin-gonic/gin"
	"mouban/common"
	"mouban/model"
)

func init() {
	err := common.DB.AutoMigrate(&model.Book{})
	if err != nil {
		panic("初始化数据库失败" + err.Error())
	}

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
