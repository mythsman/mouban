package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mouban/dao"
	"mouban/util"
	"time"
)

func runFallback() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Info(r, "fallback agent crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Hour * 1)
	}()
	cnt := dao.CasOrphanSchedule(time.Hour * 6)
	logrus.Info(cnt, "orphan schedule reset")
}

func init() {
	if !viper.GetBool("agent.enable") {
		logrus.Info("fallback agent disabled")
		return
	}

	go func() {
		for {
			runFallback()
		}
	}()

	logrus.Info("fallback agent enabled")
}
