package agent

import (
	"mouban/consts"
	"mouban/dao"
	"mouban/model"
	"mouban/util"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func userPendingSelector(ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("user pending selector panic", r, "user pending selector", "crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	schedule := dao.SearchScheduleByStatus(consts.TypeUser.Code, consts.ScheduleToCrawl.Code)
	if schedule != nil {
		logrus.Infoln("pending user found", schedule.DoubanId)
		changed := dao.CasScheduleStatus(schedule.DoubanId, consts.TypeUser.Code, consts.ScheduleCrawling.Code, *schedule.Status)
		if changed {
			ch <- schedule
		}
	} else {
		time.Sleep(10 * time.Second)
	}
}

func userRetrySelector(ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("user retry selector panic", r, "user retry selector", "crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	schedule := dao.SearchScheduleByAll(consts.TypeUser.Code, consts.ScheduleCrawled.Code, consts.ScheduleUnready.Code)
	if schedule != nil {
		logrus.Infoln("retry user found", schedule.DoubanId)
		changed := dao.CasScheduleStatus(schedule.DoubanId, consts.TypeUser.Code, consts.ScheduleCrawling.Code, *schedule.Status)
		if changed {
			ch <- schedule
		}
	} else {
		time.Sleep(time.Minute)
	}
}

func userDiscoverSelector(ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("user discover selector panic", r, "user discover selector", "crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	schedule := dao.SearchScheduleByStatus(consts.TypeUser.Code, consts.ScheduleCanCrawl.Code)

	if schedule != nil {
		logrus.Infoln("discover user found", schedule.DoubanId)
		changed := dao.CasScheduleStatus(schedule.DoubanId, consts.TypeUser.Code, consts.ScheduleCrawling.Code, *schedule.Status)
		if changed {
			ch <- schedule
		}
	} else {
		time.Sleep(time.Minute)
	}
}

func userWorker(ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("user worker panic", r, "user worker crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	for schedule := range ch {
		t := consts.ParseType(schedule.Type)
		logrus.Infoln("start process user", strconv.FormatUint(schedule.DoubanId, 10))
		processUser(schedule.DoubanId)
		dao.CasScheduleStatus(schedule.DoubanId, t.Code, consts.ScheduleCrawled.Code, consts.ScheduleCrawling.Code)
		logrus.Infoln("end process user", strconv.FormatUint(schedule.DoubanId, 10))
	}
}

func init() {
	if !viper.GetBool("agent.enable") {
		logrus.Infoln("user agent disabled")
		return
	}

	commonCh := make(chan *model.Schedule)
	discoverCh := make(chan *model.Schedule)

	go func() {
		for range time.NewTicker(time.Second).C {
			userPendingSelector(commonCh)
		}
	}()

	go func() {
		for range time.NewTicker(time.Second).C {
			userRetrySelector(commonCh)
		}
	}()

	go func() {
		for range time.NewTicker(time.Second).C {
			userWorker(commonCh)
		}
	}()

	if viper.GetBool("agent.flow.discover") {
		go func() {
			for range time.NewTicker(time.Second).C {
				userDiscoverSelector(discoverCh)
			}
		}()
		go func() {
			for range time.NewTicker(time.Second).C {
				userWorker(discoverCh)
			}
		}()
	}

	logrus.Infoln("user agent enabled")
}
