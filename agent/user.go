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
	} else {
		time.Sleep(time.Second * 10)
	}
}

func runUserStatus() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, "flow agent crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	pendingUser := dao.SearchScheduleByStatus(consts.TypeUser.Code, consts.ScheduleToCrawl.Code)
	if pendingUser == nil {
		retryUser := dao.SearchScheduleByAll(consts.TypeUser.Code, consts.ScheduleCrawled.Code, consts.ScheduleUnready.Code)
		if retryUser != nil {
			changed := dao.CasScheduleStatus(retryUser.DoubanId, retryUser.Type, consts.ScheduleToCrawl.Code, consts.ScheduleCrawled.Code)
			if changed {
				logrus.Infoln("flow retry user", retryUser.DoubanId)
			}
		} else {
			if viper.GetBool("agent.flow.discover") {
				discoverUser := dao.SearchScheduleByStatus(consts.TypeUser.Code, consts.ScheduleCanCrawl.Code)
				if discoverUser != nil {
					changed := dao.CasScheduleStatus(discoverUser.DoubanId, consts.TypeUser.Code, consts.ScheduleToCrawl.Code, consts.ScheduleCanCrawl.Code)
					if changed {
						logrus.Infoln("flow discover user", discoverUser.DoubanId)
					}
				}
			}
		}
	}
}

func init() {
	if !viper.GetBool("agent.enable") {
		logrus.Infoln("user agent disabled")
		return
	}

	go func() {
		for range time.NewTicker(time.Second).C {
			runUser()
		}
	}()

	go func() {
		for range time.NewTicker(time.Second * 5).C {
			runUserStatus()
		}
	}()

	logrus.Infoln("user agent enabled")
}
