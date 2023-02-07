package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
		time.Sleep(time.Second * 1)
	}()

	pendingBook := dao.SearchScheduleByStatus(consts.TypeBook.Code, consts.ScheduleStatusToCrawl)
	if pendingBook == nil {
		retryBook := dao.SearchScheduleByAll(consts.TypeBook.Code, consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
		if retryBook != nil {
			changed := dao.CasScheduleStatus(retryBook.DoubanId, retryBook.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
			if changed {
				logrus.Infoln("flow retry book ", retryBook.DoubanId)
			}
		} else {
			discoverBook := dao.SearchScheduleByStatus(consts.TypeBook.Code, consts.ScheduleStatusCanCrawl)
			if discoverBook != nil {
				changed := dao.CasScheduleStatus(discoverBook.DoubanId, discoverBook.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCanCrawl)
				if changed {
					logrus.Infoln("flow discover book", discoverBook.DoubanId)
				}
			}
		}
	}

	pendingMovie := dao.SearchScheduleByStatus(consts.TypeMovie.Code, consts.ScheduleStatusToCrawl)
	if pendingMovie == nil {
		retryMovie := dao.SearchScheduleByAll(consts.TypeMovie.Code, consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
		if retryMovie != nil {
			changed := dao.CasScheduleStatus(retryMovie.DoubanId, retryMovie.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
			if changed {
				logrus.Infoln("flow retry movie ", retryMovie.DoubanId)
			}
		} else {
			discoverMovie := dao.SearchScheduleByStatus(consts.TypeMovie.Code, consts.ScheduleStatusCanCrawl)
			if discoverMovie != nil {
				changed := dao.CasScheduleStatus(discoverMovie.DoubanId, discoverMovie.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCanCrawl)
				if changed {
					logrus.Infoln("flow discover movie", discoverMovie.DoubanId)
				}
			}
		}
	}

	pendingGame := dao.SearchScheduleByStatus(consts.TypeGame.Code, consts.ScheduleStatusToCrawl)
	if pendingGame == nil {
		retryGame := dao.SearchScheduleByAll(consts.TypeGame.Code, consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
		if retryGame != nil {
			changed := dao.CasScheduleStatus(retryGame.DoubanId, retryGame.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
			if changed {
				logrus.Infoln("flow retry game ", retryGame.DoubanId)
			}
		} else {
			discoverGame := dao.SearchScheduleByStatus(consts.TypeGame.Code, consts.ScheduleStatusCanCrawl)
			if discoverGame != nil {
				changed := dao.CasScheduleStatus(discoverGame.DoubanId, discoverGame.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCanCrawl)
				if changed {
					logrus.Infoln("flow discover game", discoverGame.DoubanId)
				}
			}
		}
	}

	pendingSong := dao.SearchScheduleByStatus(consts.TypeSong.Code, consts.ScheduleStatusToCrawl)
	if pendingSong == nil {
		retrySong := dao.SearchScheduleByAll(consts.TypeSong.Code, consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
		if retrySong != nil {
			changed := dao.CasScheduleStatus(retrySong.DoubanId, retrySong.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
			if changed {
				logrus.Infoln("flow retry song ", retrySong.DoubanId)
			}
		} else {
			discoverSong := dao.SearchScheduleByStatus(consts.TypeSong.Code, consts.ScheduleStatusCanCrawl)
			if discoverSong != nil {
				changed := dao.CasScheduleStatus(discoverSong.DoubanId, discoverSong.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCanCrawl)
				if changed {
					logrus.Infoln("flow discover song", discoverSong.DoubanId)
				}
			}
		}
	}

	pendingUser := dao.SearchScheduleByStatus(consts.TypeUser.Code, consts.ScheduleStatusToCrawl)
	if pendingUser == nil {
		retryUser := dao.SearchScheduleByAll(consts.TypeUser.Code, consts.ScheduleStatusCrawled, consts.ScheduleResultUnready)
		if retryUser != nil {
			changed := dao.CasScheduleStatus(retryUser.DoubanId, retryUser.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
			if changed {
				logrus.Infoln("flow retry user ", retryUser.DoubanId)
			}
		} else {
			if viper.GetBool("agent.flow.discover") {
				discoverUser := dao.SearchScheduleByStatus(consts.TypeUser.Code, consts.ScheduleStatusCanCrawl)
				if discoverUser != nil {
					changed := dao.CasScheduleStatus(discoverUser.DoubanId, consts.TypeUser.Code, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCanCrawl)
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
		for {
			runFlow()
		}
	}()

	logrus.Infoln("flow agent enabled")
}
