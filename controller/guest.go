package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ListUserMusic(ctx *gin.Context) {
	id := ctx.Query("id")
	fmt.Println("id: ", id)
}

func ListUserMovie(ctx *gin.Context) {
	id := ctx.Query("id")
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
