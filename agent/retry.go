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
			log.Println(r, "retry agent crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Second * 1)
	}()

	pendingBook := dao.SearchScheduleByStatus(consts.TypeBook.Code, consts.ScheduleStatusToCrawl)
	if pendingBook == nil {
		retryBook := dao.SearchScheduleByAll(consts.TypeBook.Code, consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
		if retryBook != nil {
			changed := dao.CasScheduleStatus(retryBook.DoubanId, retryBook.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
			if changed {
				log.Println("process retry book ", retryBook.DoubanId)
			}
		}
	}

	pendingMovie := dao.SearchScheduleByStatus(consts.TypeMovie.Code, consts.ScheduleStatusToCrawl)
	if pendingMovie == nil {
		retryMovie := dao.SearchScheduleByAll(consts.TypeMovie.Code, consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
		if retryMovie != nil {
			changed := dao.CasScheduleStatus(retryMovie.DoubanId, retryMovie.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
			if changed {
				log.Println("process retry movie ", retryMovie.DoubanId)
			}
		}
	}

	pendingGame := dao.SearchScheduleByStatus(consts.TypeGame.Code, consts.ScheduleStatusToCrawl)
	if pendingGame == nil {
		retryGame := dao.SearchScheduleByAll(consts.TypeGame.Code, consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
		if retryGame != nil {
			changed := dao.CasScheduleStatus(retryGame.DoubanId, retryGame.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
			if changed {
				log.Println("process retry game ", retryGame.DoubanId)
			}
		}
	}

	pendingUser := dao.SearchScheduleByStatus(consts.TypeUser.Code, consts.ScheduleStatusToCrawl)
	if pendingUser == nil {
		retryUser := dao.SearchScheduleByAll(consts.TypeUser.Code, consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
		if retryUser != nil {
			changed := dao.CasScheduleStatus(retryUser.DoubanId, retryUser.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
			if changed {
				log.Println("process retry user ", retryUser.DoubanId)
			}
		} else {
			if viper.GetString("agent.discover") == "true" {
				discoverUser := dao.SearchScheduleByStatus(consts.TypeUser.Code, consts.ScheduleStatusCanCrawl)
				if discoverUser != nil {
					changed := dao.CasScheduleStatus(discoverUser.DoubanId, consts.TypeUser.Code, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCanCrawl)
					if changed {
						log.Println("process discover user ", discoverUser.DoubanId)
					}
				}
			}
		}
	}

}
func init() {
	if viper.GetString("agent.enable") != "true" {
		log.Println("retry agent disabled")
		return
	}
	go func() {
		for {
			runRetry()
		}
	}()

	log.Println("retry agent enabled")
}
