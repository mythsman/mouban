package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"mouban/consts"
	"mouban/dao"
	"mouban/util"
	"time"
)

func runFlow() {
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
		logrus.Infoln("flow agent disabled")
		return
	}
	go func() {
		for range time.NewTicker(time.Second * 5).C {
			runFlow()
		}
	}()

	logrus.Infoln("flow agent enabled")
}
