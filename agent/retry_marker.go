package agent

import (
	"github.com/spf13/viper"
	"log"
	"mouban/consts"
	"mouban/dao"
	"mouban/util"
	"time"
)

func runRetry() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, "retry marker crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Second * 5)
	}()

	pendingBook := dao.SearchScheduleByStatus(consts.TypeBook, consts.ScheduleStatusToCrawl)
	if pendingBook == nil {
		retryBook := dao.SearchScheduleByAll(consts.TypeBook, consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
		if retryBook != nil {
			changed := dao.CasScheduleStatus(retryBook.DoubanId, retryBook.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
			if changed {
				log.Printf("process retry book %d %d\n", retryBook.Type, retryBook.DoubanId)
			}
		}
	}

	pendingMovie := dao.SearchScheduleByStatus(consts.TypeMovie, consts.ScheduleStatusToCrawl)
	if pendingMovie == nil {
		retryMovie := dao.SearchScheduleByAll(consts.TypeMovie, consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
		if retryMovie != nil {
			changed := dao.CasScheduleStatus(retryMovie.DoubanId, retryMovie.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
			if changed {
				log.Printf("process retry movie %d %d\n", retryMovie.Type, retryMovie.DoubanId)
			}
		}
	}
	pendingGame := dao.SearchScheduleByStatus(consts.TypeGame, consts.ScheduleStatusToCrawl)
	if pendingGame == nil {
		retryGame := dao.SearchScheduleByAll(consts.TypeGame, consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
		if retryGame != nil {
			changed := dao.CasScheduleStatus(retryGame.DoubanId, retryGame.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
			if changed {
				log.Printf("process retry game %d %d\n", retryGame.Type, retryGame.DoubanId)
			}
		}
	}

	pendingUser := dao.SearchScheduleByStatus(consts.TypeUser, consts.ScheduleStatusToCrawl)
	if pendingUser == nil {
		retryUser := dao.SearchScheduleByAll(consts.TypeUser, consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
		if retryUser != nil {
			changed := dao.CasScheduleStatus(retryUser.DoubanId, retryUser.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
			if changed {
				log.Printf("process retry user %d %d\n", retryUser.Type, retryUser.DoubanId)
			}
		} else {
			discoverUser := dao.SearchScheduleByStatus(consts.TypeUser, consts.ScheduleStatusCanCrawl)
			if discoverUser != nil {
				changed := dao.CasScheduleStatus(discoverUser.DoubanId, consts.TypeUser, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCanCrawl)
				if changed {
					log.Printf("process discover user %d\n", discoverUser.DoubanId)
				}
			}
		}
	}

}
func init() {
	if viper.GetString("agent.enable") != "true" {
		log.Println("retry marker disabled")
		return
	}
	go func() {
		for {
			runRetry()
		}
	}()

	log.Println("retry marker enabled")
}
