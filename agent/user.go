package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mouban/consts"
	"mouban/dao"
	"mouban/util"
	"strconv"
	"time"
)

func runUser() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, "user agent crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Second * 1)
	}()
	schedule := dao.SearchScheduleByAll(consts.TypeUser.Code, consts.ScheduleToCrawl.Code, consts.ScheduleReady.Code)
	if schedule == nil {
		schedule = dao.SearchScheduleByAll(consts.TypeUser.Code, consts.ScheduleToCrawl.Code, consts.ScheduleUnready.Code)
	}
	if schedule != nil {
		changed := dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleCrawling.Code, consts.ScheduleToCrawl.Code)
		if changed {
			logrus.Infoln("start process user", strconv.FormatUint(schedule.DoubanId, 10))
			processUser(schedule.DoubanId)
			dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleCrawled.Code, consts.ScheduleCrawling.Code)
			logrus.Infoln("end process user", strconv.FormatUint(schedule.DoubanId, 10))
		}
	}
}

func init() {
	if !viper.GetBool("agent.enable") {
		logrus.Infoln("user agent disabled")
		return
	}

	go func() {
		for {
			runUser()
		}
	}()

	logrus.Infoln("user agent enabled")
}
