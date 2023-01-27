package agent

import (
	"github.com/spf13/viper"
	"log"
	"mouban/dao"
	"mouban/util"
	"time"
)

func runFallback() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, "fallback agent crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Hour * 1)
	}()
	cnt := dao.CasOrphanSchedule(time.Hour * 6)
	log.Println(cnt, "orphan schedule reset")
}

func init() {
	if !viper.GetBool("agent.enable") {
		log.Println("fallback agent disabled")
		return
	}

	go func() {
		for {
			runFallback()
		}
	}()

	log.Println("fallback agent enabled")
}
