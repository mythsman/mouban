package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mouban/consts"
	"mouban/dao"
	"mouban/model"
	"mouban/util"
	"strconv"
	"time"
)

func userPendingSelector(ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, "user pending selector", "crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	schedule := dao.SearchScheduleByStatus(consts.TypeUser.Code, consts.ScheduleToCrawl.Code)
	if schedule != nil {
		ch <- schedule
	} else {
		time.Sleep(10 * time.Second)
	}
}

func userRetrySelector(ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, "user retry selector", "crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	schedule := dao.SearchScheduleByAll(consts.TypeUser.Code, consts.ScheduleCrawled.Code, consts.ScheduleUnready.Code)
	if schedule != nil {
		ch <- schedule
	} else {
		time.Sleep(time.Minute)
	}
}

func userDiscoverSelector(ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, "user discover selector", "crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	schedule := dao.SearchScheduleByStatus(consts.TypeUser.Code, consts.ScheduleCanCrawl.Code)

	if schedule != nil {
		ch <- schedule
	} else {
		time.Sleep(time.Minute)
	}
}

func userWorker(ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, "user worker crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	for schedule := range ch {
		t := consts.ParseType(schedule.Type)
		changed := dao.CasScheduleStatus(schedule.DoubanId, t.Code, consts.ScheduleCrawling.Code, consts.ScheduleToCrawl.Code)
		if changed {
			logrus.Infoln("start process user", strconv.FormatUint(schedule.DoubanId, 10))
			processUser(schedule.DoubanId)
			dao.CasScheduleStatus(schedule.DoubanId, t.Code, consts.ScheduleCrawled.Code, consts.ScheduleCrawling.Code)
			logrus.Infoln("end process user", strconv.FormatUint(schedule.DoubanId, 10))
		}
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
			userWorker(discoverCh)
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
				userWorker(commonCh)
			}
		}()
	}

	logrus.Infoln("user agent enabled")
}
