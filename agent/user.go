package agent

import (
	"github.com/spf13/viper"
	"mouban/consts"
	"mouban/dao"
	"mouban/log"
	"mouban/util"
	"strconv"
	"time"
)

func runUser() {
	defer func() {
		if r := recover(); r != nil {
			log.Info(r, "user agent crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Second * 1)
	}()
	schedule := dao.SearchScheduleByAll(consts.TypeUser.Code, consts.ScheduleStatusToCrawl, consts.ScheduleResultReady)
	if schedule == nil {
		schedule = dao.SearchScheduleByAll(consts.TypeUser.Code, consts.ScheduleStatusToCrawl, consts.ScheduleResultUnready)
	}
	if schedule != nil {
		changed := dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawling, consts.ScheduleStatusToCrawl)
		if changed {
			log.Info("start process user" + strconv.FormatUint(schedule.DoubanId, 10))
			processUser(schedule.DoubanId, schedule.Result == consts.ScheduleResultUnready)
			dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawled, consts.ScheduleStatusCrawling)
			log.Info("end process user" + strconv.FormatUint(schedule.DoubanId, 10))
		}
	}
}

func init() {
	if !viper.GetBool("agent.enable") {
		log.Info("user agent disabled")
		return
	}

	go func() {
		for {
			runUser()
		}
	}()

	log.Info("user agent enabled")
}
