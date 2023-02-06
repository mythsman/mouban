package agent

import (
	"github.com/spf13/viper"
	"mouban/dao"
	"mouban/log"
	"mouban/util"
	"time"
)

func runFallback() {
	defer func() {
		if r := recover(); r != nil {
			log.Info(r, "fallback agent crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Hour * 1)
	}()
	cnt := dao.CasOrphanSchedule(time.Hour * 6)
	log.Info(cnt, "orphan schedule reset")
}

func init() {
	if !viper.GetBool("agent.enable") {
		log.Info("fallback agent disabled")
		return
	}

	go func() {
		for {
			runFallback()
		}
	}()

	log.Info("fallback agent enabled")
}
