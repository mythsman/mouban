package agent

import (
	"github.com/spf13/viper"
	"log"
	"mouban/consts"
	"mouban/dao"
	"mouban/util"
	"time"
)

func runRetry() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, "retry marker crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Second * 5)
	}()
	schedule := dao.SearchSchedule(consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
	if schedule != nil {
		changed := dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
		if changed {
			log.Printf("process retry %d %d \n", schedule.Type, schedule.DoubanId)
		}
	}
}
func init() {
	if viper.GetString("agent.enable") != "true" {
		log.Println("retry marker disabled")
		return
	}
	go func() {
		for {
			runRetry()
		}
	}()

	log.Println("retry marker enabled")
}
