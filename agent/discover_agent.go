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

func runDiscover() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, "discover agent crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Second * 60)
	}()
	schedule := dao.SearchScheduleByStatus(consts.TypeUser, consts.ScheduleStatusCanCrawl)
	if schedule != nil {
		changed := dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawling, consts.ScheduleStatusCanCrawl)
		if changed {
			log.Println("start process discover " + strconv.FormatUint(schedule.DoubanId, 10))
			processUser(schedule.DoubanId)
			dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawled, consts.ScheduleStatusCrawling)
			log.Println("end process discover " + strconv.FormatUint(schedule.DoubanId, 10))
		}
	}
}
func init() {
	if viper.GetString("agent.enable") != "true" {
		log.Println("discover agent disabled")
		return
	}
	go func() {
		for {
			runDiscover()
		}
	}()

	log.Println("discover agent enabled")
}
