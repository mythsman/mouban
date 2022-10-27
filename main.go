package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mouban/common"
)

func main() {

	router := gin.Default()

	router = common.CollectRoute(router)

	panic(router.Run(":" + viper.GetString("server.port")))

}
