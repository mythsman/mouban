package common

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func initFromYaml() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir)
	viper.AddConfigPath(filepath.Join(workDir, "../"))
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Infoln("viper init error")
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

func InitConfig() {
	initFromYaml()
	initFromEnv()
}
