package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	_ "mouban/common"
)

func main() {

	router := gin.Default()

	router = CollectRoute(router)

	panic(router.Run(":" + viper.GetString("server.port")))

}
