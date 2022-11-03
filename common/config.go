package common

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir)
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("获取配置文件错误")
		panic(err)
	}
	log.Println("配置初始化成功")

}
