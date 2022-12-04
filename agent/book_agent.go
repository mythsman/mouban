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

func runBook() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, "book agent crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Second * 5)
	}()

	schedule := dao.SearchScheduleByStatus(consts.TypeBook, consts.ScheduleStatusToCrawl)
	if schedule != nil {
		changed := dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawling, consts.ScheduleStatusToCrawl)
		if changed {
			log.Println("start process book " + strconv.FormatUint(schedule.DoubanId, 10))
			processBook(schedule.DoubanId)
			dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawled, consts.ScheduleStatusCrawling)
			log.Println("end process book " + strconv.FormatUint(schedule.DoubanId, 10))
		}
	}
}
func init() {
	if viper.GetString("agent.enable") != "true" {
		log.Println("book agent disabled")
		return
	}
	go func() {
		for {
			runBook()
		}
	}()

	log.Println("book agent enabled")
}
