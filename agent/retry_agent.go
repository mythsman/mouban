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

func runRetry() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, "retry agent crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Second * 5)
	}()
	schedule := dao.SearchSchedule(consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
	if schedule != nil {
		changed := dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawling, consts.ScheduleStatusCanCrawl)
		if changed {
			log.Println("start process retry " + strconv.FormatUint(schedule.DoubanId, 10))
			switch schedule.Type {
			case consts.TypeUser:
				processUser(schedule.DoubanId)
				break
			case consts.TypeBook:
				processBook(schedule.DoubanId)
				break
			case consts.TypeGame:
				processGame(schedule.DoubanId)
				break
			case consts.TypeMovie:
				processMovie(schedule.DoubanId)
				break
			}
		}
	}
}
func init() {
	if viper.GetString("agent.enable") != "true" {
		log.Println("retry agent disabled")
		return
	}
	go func() {
		for {
			runRetry()
		}
	}()

	log.Println("retry agent enabled")
}