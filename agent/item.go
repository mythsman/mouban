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

func itemPendingSelector(t consts.Type, ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, "item pending selector for", t.Name, "crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	schedule := dao.SearchScheduleByStatus(t.Code, consts.ScheduleToCrawl.Code)
	if schedule != nil {
		logrus.Infoln("pending", t.Name, "item found", schedule.DoubanId)
		changed := dao.CasScheduleStatus(schedule.DoubanId, t.Code, consts.ScheduleCrawling.Code, *schedule.Status)
		if changed {
			ch <- schedule
		}
	} else {
		logrus.Infoln("item", t.Name, "pending selector idle")
		time.Sleep(10 * time.Second)
	}
}

func itemRetrySelector(t consts.Type, ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, "item retry selector for", t.Name, "crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	schedule := dao.SearchScheduleByAll(t.Code, consts.ScheduleCrawled.Code, consts.ScheduleUnready.Code)
	if schedule != nil {
		logrus.Infoln("retry", t.Name, "item found", schedule.DoubanId)
		changed := dao.CasScheduleStatus(schedule.DoubanId, t.Code, consts.ScheduleCrawling.Code, *schedule.Status)
		if changed {
			ch <- schedule
		}
	} else {
		logrus.Infoln("item", t.Name, "retry selector idle")
		time.Sleep(time.Minute)
	}
}

func itemDiscoverSelector(t consts.Type, ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, "item discover selector for", t.Name, "crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	schedule := dao.SearchScheduleByStatus(t.Code, consts.ScheduleCanCrawl.Code)

	if schedule != nil {
		logrus.Infoln("discover", t.Name, "item found", schedule.DoubanId)
		changed := dao.CasScheduleStatus(schedule.DoubanId, t.Code, consts.ScheduleCrawling.Code, *schedule.Status)
		if changed {
			ch <- schedule
		}
	} else {
		logrus.Infoln("item", t.Name, "discover selector idle")
		time.Sleep(time.Minute)
	}
}

func itemWorker(index int, ch chan *model.Schedule) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, "item worker (", index, ") crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	for schedule := range ch {
		t := consts.ParseType(schedule.Type)
		logrus.Infoln("item thread", index, "start", t.Name, strconv.FormatUint(schedule.DoubanId, 10))
		processItem(schedule.Type, schedule.DoubanId)
		dao.CasScheduleStatus(schedule.DoubanId, t.Code, consts.ScheduleCrawled.Code, consts.ScheduleCrawling.Code)
		logrus.Infoln("item thread", index, "end", t.Name, strconv.FormatUint(schedule.DoubanId, 10))
	}
}

func init() {
	if !viper.GetBool("agent.enable") {
		logrus.Infoln("item agent disabled")
		return
	}

	ch := make(chan *model.Schedule)

	types := []consts.Type{consts.TypeBook, consts.TypeMovie, consts.TypeGame, consts.TypeSong}
	for i := range types {
		i := i
		go func() {
			for range time.NewTicker(time.Second).C {
				itemPendingSelector(types[i], ch)
			}
		}()
		go func() {
			for range time.NewTicker(time.Second).C {
				itemRetrySelector(types[i], ch)
			}
		}()
		go func() {
			for range time.NewTicker(time.Second).C {
				itemDiscoverSelector(types[i], ch)
			}
		}()
	}

	concurrency := viper.GetInt("agent.item.concurrency")
	for i := 0; i < concurrency; i++ {
		j := i + 1
		go func() {
			for range time.NewTicker(time.Second).C {
				itemWorker(j, ch)
			}
		}()
	}

	logrus.Infoln(concurrency, "item agent(s) enabled")
}
