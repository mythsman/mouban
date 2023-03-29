package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"mouban/consts"
	"mouban/dao"
	"mouban/util"
	"strconv"
	"time"
)

func runItemAgent(index int) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, "item agent", index, "crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	types := []consts.Type{consts.TypeBook, consts.TypeMovie, consts.TypeGame, consts.TypeSong}
	rand.Shuffle(len(types), func(i, j int) { types[i], types[j] = types[j], types[i] })

	found := false
	for _, t := range types {
		schedule := dao.SearchScheduleByStatus(t.Code, consts.ScheduleToCrawl.Code)
		if schedule != nil {
			changed := dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleCrawling.Code, consts.ScheduleToCrawl.Code)
			if changed {
				found = true
				logrus.Infoln("item thread", index, "start", t.Name, strconv.FormatUint(schedule.DoubanId, 10))
				processItem(t.Code, schedule.DoubanId)
				dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleCrawled.Code, consts.ScheduleCrawling.Code)
				logrus.Infoln("item thread", index, "end", t.Name, strconv.FormatUint(schedule.DoubanId, 10))
			}
			break
		}
	}
	if !found {
		time.Sleep(time.Second * 10)
	}
}

func runItemStatus() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, "flow agent crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	types := []consts.Type{consts.TypeBook, consts.TypeMovie, consts.TypeGame, consts.TypeSong}
	rand.Shuffle(len(types), func(i, j int) { types[i], types[j] = types[j], types[i] })

	for _, t := range types {
		pendingItem := dao.SearchScheduleByStatus(t.Code, consts.ScheduleToCrawl.Code)
		if pendingItem == nil {
			retryItem := dao.SearchScheduleByAll(t.Code, consts.ScheduleCrawled.Code, consts.ScheduleUnready.Code)
			if retryItem != nil {
				changed := dao.CasScheduleStatus(retryItem.DoubanId, retryItem.Type, consts.ScheduleToCrawl.Code, consts.ScheduleCrawled.Code)
				if changed {
					logrus.Infoln("flow retry", t.Name, retryItem.DoubanId)
				}
			} else {
				concurrency := viper.GetInt("agent.item.concurrency")
				for i := 0; i < concurrency; i++ {
					discoverItem := dao.SearchScheduleByStatus(t.Code, consts.ScheduleCanCrawl.Code)
					if discoverItem != nil {
						changed := dao.CasScheduleStatus(discoverItem.DoubanId, discoverItem.Type, consts.ScheduleToCrawl.Code, consts.ScheduleCanCrawl.Code)
						if changed {
							logrus.Infoln("flow discover", t.Name, discoverItem.DoubanId)
						}
					} else {
						break
					}
				}
			}
		}
	}
}

func init() {
	if !viper.GetBool("agent.enable") {
		logrus.Infoln("item agent disabled")
		return
	}

	concurrency := viper.GetInt("agent.item.concurrency")
	for i := 0; i < concurrency; i++ {
		j := i + 1
		go func() {
			for range time.NewTicker(time.Second).C {
				runItemAgent(j)
			}
		}()
	}

	go func() {
		for range time.NewTicker(time.Second * 5).C {
			runItemStatus()
		}
	}()

	logrus.Infoln(concurrency, "item agent(s) enabled")
}
