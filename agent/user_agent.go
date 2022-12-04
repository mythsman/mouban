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

func runUser() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, "user agent crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Second * 5)
	}()
	schedule := dao.SearchScheduleByStatus(consts.TypeUser, consts.ScheduleStatusToCrawl)
	if schedule != nil {
		changed := dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawling, consts.ScheduleStatusToCrawl)
		if changed {
			log.Println("start process user " + strconv.FormatUint(schedule.DoubanId, 10))
			processUser(schedule.DoubanId, schedule.Result == consts.ScheduleResultUnready)
			dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawled, consts.ScheduleStatusCrawling)
			log.Println("end process user " + strconv.FormatUint(schedule.DoubanId, 10))
		}
	}
}
func init() {
	if viper.GetString("agent.enable") != "true" {
		log.Println("user agent disabled")
		return
	}

	go func() {
		for {
			runUser()
		}
	}()

	log.Println("user agent enabled")
}
