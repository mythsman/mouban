package common

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

func initFromYaml() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir)
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("viper init error")
		panic(err)
	}
}

func initFromEnv() {
	for _, s := range viper.AllKeys() {
		yamlCfg := strings.ReplaceAll(s, ".", "__")
		envCfg := os.Getenv(yamlCfg)
		if envCfg != "" {
			viper.Set(s, envCfg)
		}
	}
}

func init() {
	initFromYaml()
	initFromEnv()
	log.Println("config init success")
}
