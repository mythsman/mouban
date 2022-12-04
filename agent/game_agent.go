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

func runGame() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, "game agent crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()
	schedule := dao.SearchScheduleByStatus(consts.TypeGame, consts.ScheduleStatusToCrawl)
	if schedule == nil {
		time.Sleep(time.Second * 5)
	} else {
		changed := dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawling, consts.ScheduleStatusToCrawl)
		if changed {
			log.Println("start process game " + strconv.FormatUint(schedule.DoubanId, 10))
			processBook(schedule.DoubanId)
		}
	}
}
func init() {
	if viper.GetString("agent.enable") != "true" {
		log.Println("game agent disabled")
		return
	}
	go func() {
		for {
			runGame()
		}
	}()

	log.Println("game agent enabled")
}
