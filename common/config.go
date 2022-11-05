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
		log.Println("viper init error")
		panic(err)
	}
	log.Println("viper init success")

}
