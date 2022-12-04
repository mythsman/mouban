package agent

import (
	"github.com/spf13/viper"
	"log"
	"mouban/consts"
	"mouban/dao"
	"mouban/util"
	"strconv"
	"time"
)

func runMovie() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, "movie agent crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Second * 5)
	}()
	schedule := dao.SearchScheduleByStatus(consts.TypeMovie, consts.ScheduleStatusToCrawl)
	if schedule != nil {
		changed := dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawling, consts.ScheduleStatusToCrawl)
		if changed {
			log.Println("start process movie " + strconv.FormatUint(schedule.DoubanId, 10))
			processMovie(schedule.DoubanId)
		}
	}
}
func init() {
	if viper.GetString("agent.enable") != "true" {
		log.Println("movie agent disabled")
		return
	}
	go func() {
		for {
			runMovie()
		}
	}()

	log.Println("movie agent enabled")
}
