package common

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func init() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("获取配置文件错误")
		panic(err)
	}
	fmt.Println("配置初始化成功")

}
