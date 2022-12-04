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

func runDiscover() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, "discover marker crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Second * 5)
	}()

	currentUser := dao.SearchScheduleByStatus(consts.TypeUser, consts.ScheduleStatusCrawling)
	if currentUser == nil {
		pendingUser := dao.SearchScheduleByStatus(consts.TypeUser, consts.ScheduleStatusToCrawl)
		if pendingUser == nil {
			nextUser := dao.SearchScheduleByStatus(consts.TypeUser, consts.ScheduleStatusCanCrawl)
			if nextUser != nil {
				changed := dao.CasScheduleStatus(nextUser.DoubanId, consts.TypeUser, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCanCrawl)
				if changed {
					log.Println("process discover " + strconv.FormatUint(nextUser.DoubanId, 10))
				}
			}
		}
	}

}
func init() {
	if viper.GetString("agent.enable") != "true" {
		log.Println("discover agent disabled")
		return
	}
	go func() {
		for {
			runDiscover()
		}
	}()

	log.Println("discover agent enabled")
}
