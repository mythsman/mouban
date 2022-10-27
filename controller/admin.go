package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func GetOverview(ctx *gin.Context) {

}

func CrawlUser(ctx *gin.Context) {
	id := ctx.Query("id")
	fmt.Println("id: ", id)
}

func CrawlGame(ctx *gin.Context) {
	id := ctx.Query("id")
	fmt.Println("id: ", id)
}

func CrawlBook(ctx *gin.Context) {
	id := ctx.Query("id")
	fmt.Println("id: ", id)
}

func CrawlMovie(ctx *gin.Context) {
	id := ctx.Query("id")
	fmt.Println("id: ", id)
}

func CrawlMusic(ctx *gin.Context) {
	id := ctx.Query("id")
	fmt.Println("id: ", id)
}
