package common

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

func initFromYaml() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir)
	viper.AddConfigPath(filepath.Join(workDir, "../"))
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Info("viper init error")
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
