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
			logrus.Errorln(r, "fallback agent crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Hour * 1)
	}()
	cnt := dao.CasOrphanSchedule(time.Hour * 6)
	logrus.Infoln(cnt, "orphan schedule reset")
}

func init() {
	if !viper.GetBool("agent.enable") {
		logrus.Infoln("fallback agent disabled")
		return
	}

	go func() {
		for {
			runFallback()
		}
	}()

	logrus.Infoln("fallback agent enabled")
}
