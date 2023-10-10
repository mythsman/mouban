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

func userSelector(ch chan *model.Schedule, done chan bool) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("user selector panic", r, util.GetCurrentGoroutineStack())
		}
	}()

	<-done

	for {
		pendingSchedule := dao.SearchScheduleByStatus(consts.TypeUser.Code, consts.ScheduleToCrawl.Code)
		if pendingSchedule != nil {
			logrus.Infoln("pending user found", pendingSchedule.DoubanId)
			changed := dao.CasScheduleStatus(pendingSchedule.DoubanId, consts.TypeUser.Code, consts.ScheduleCrawling.Code, *pendingSchedule.Status)
			if changed {
				ch <- pendingSchedule
				return
			}
		}

		retrySchedule := dao.SearchScheduleByAll(consts.TypeUser.Code, consts.ScheduleCrawled.Code, consts.ScheduleUnready.Code)

		if retrySchedule != nil {
			logrus.Infoln("retry user found", retrySchedule.DoubanId)
			changed := dao.CasScheduleStatus(retrySchedule.DoubanId, consts.TypeUser.Code, consts.ScheduleCrawling.Code, *retrySchedule.Status)
			if changed {
				ch <- retrySchedule
				return
			}
		}

		if viper.GetBool("agent.flow.discover") {
			discoverSchedule := dao.SearchScheduleByStatus(consts.TypeUser.Code, consts.ScheduleCanCrawl.Code)
			if discoverSchedule != nil {
				logrus.Infoln("discover user found", discoverSchedule.DoubanId)
				changed := dao.CasScheduleStatus(discoverSchedule.DoubanId, consts.TypeUser.Code, consts.ScheduleCrawling.Code, *discoverSchedule.Status)
				if changed {
					ch <- discoverSchedule
					return
				}
			}
		}

		time.Sleep(time.Minute)
	}
}

func userWorker(index int, ch chan *model.Schedule, done chan bool) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("user worker panic", r, "user worker (", index, ") crashed  => ", util.GetCurrentGoroutineStack())
		}
		done <- true
	}()

	schedule := <-ch

	t := consts.ParseType(schedule.Type)
	logrus.Infoln("user thread", index, "start", strconv.FormatUint(schedule.DoubanId, 10))
	processUser(schedule.DoubanId)
	dao.CasScheduleStatus(schedule.DoubanId, t.Code, consts.ScheduleCrawled.Code, consts.ScheduleCrawling.Code)
	logrus.Infoln("user thread", index, "end", strconv.FormatUint(schedule.DoubanId, 10))
}

func init() {
	if !viper.GetBool("agent.enable") {
		logrus.Infoln("user agent disabled")
		return
	}

	concurrency := viper.GetInt("agent.user.concurrency")

	ch := make(chan *model.Schedule, concurrency)
	done := make(chan bool, concurrency)

	go func() {
		for {
			userSelector(ch, done)
		}
	}()

	for i := 0; i < concurrency; i++ {
		j := i + 1
		go func() {
			for {
				userWorker(j, ch, done)
			}
		}()
		done <- true
	}

	logrus.Infoln(concurrency, "user agent(s) enabled")

}
