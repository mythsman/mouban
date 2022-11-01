package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mouban/dao"
)

func ListUserMovie(ctx *gin.Context) {
	id := ctx.Query("id")
	user := dao.GetUser(id)
	if user == nil {

	}

	fmt.Println("id: ", id)
}

func ListUserBook(ctx *gin.Context) {
	id := ctx.Query("id")
	fmt.Println("id: ", id)

}

func ListUserGame(ctx *gin.Context) {
	id := ctx.Query("id")
	fmt.Println("id: ", id)
}

