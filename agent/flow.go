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

	pendingBook := dao.SearchScheduleByStatus(consts.TypeBook.Code, consts.ScheduleToCrawl.Code)
	if pendingBook == nil {
		retryBook := dao.SearchScheduleByAll(consts.TypeBook.Code, consts.ScheduleCrawled.Code, consts.ScheduleUnready.Code)
		if retryBook != nil {
			changed := dao.CasScheduleStatus(retryBook.DoubanId, retryBook.Type, consts.ScheduleToCrawl.Code, consts.ScheduleCrawled.Code)
			if changed {
				logrus.Infoln("flow retry book ", retryBook.DoubanId)
			}
		} else {
			discoverBook := dao.SearchScheduleByStatus(consts.TypeBook.Code, consts.ScheduleCanCrawl.Code)
			if discoverBook != nil {
				changed := dao.CasScheduleStatus(discoverBook.DoubanId, discoverBook.Type, consts.ScheduleToCrawl.Code, consts.ScheduleCanCrawl.Code)
				if changed {
					logrus.Infoln("flow discover book", discoverBook.DoubanId)
				}
			}
		}
	}

	pendingMovie := dao.SearchScheduleByStatus(consts.TypeMovie.Code, consts.ScheduleToCrawl.Code)
	if pendingMovie == nil {
		retryMovie := dao.SearchScheduleByAll(consts.TypeMovie.Code, consts.ScheduleCrawled.Code, consts.ScheduleUnready.Code)
		if retryMovie != nil {
			changed := dao.CasScheduleStatus(retryMovie.DoubanId, retryMovie.Type, consts.ScheduleToCrawl.Code, consts.ScheduleCrawled.Code)
			if changed {
				logrus.Infoln("flow retry movie ", retryMovie.DoubanId)
			}
		} else {
			discoverMovie := dao.SearchScheduleByStatus(consts.TypeMovie.Code, consts.ScheduleCanCrawl.Code)
			if discoverMovie != nil {
				changed := dao.CasScheduleStatus(discoverMovie.DoubanId, discoverMovie.Type, consts.ScheduleToCrawl.Code, consts.ScheduleCanCrawl.Code)
				if changed {
					logrus.Infoln("flow discover movie", discoverMovie.DoubanId)
				}
			}
		}
	}

	pendingGame := dao.SearchScheduleByStatus(consts.TypeGame.Code, consts.ScheduleToCrawl.Code)
	if pendingGame == nil {
		retryGame := dao.SearchScheduleByAll(consts.TypeGame.Code, consts.ScheduleCrawled.Code, consts.ScheduleUnready.Code)
		if retryGame != nil {
			changed := dao.CasScheduleStatus(retryGame.DoubanId, retryGame.Type, consts.ScheduleToCrawl.Code, consts.ScheduleCrawled.Code)
			if changed {
				logrus.Infoln("flow retry game ", retryGame.DoubanId)
			}
		} else {
			discoverGame := dao.SearchScheduleByStatus(consts.TypeGame.Code, consts.ScheduleCanCrawl.Code)
			if discoverGame != nil {
				changed := dao.CasScheduleStatus(discoverGame.DoubanId, discoverGame.Type, consts.ScheduleToCrawl.Code, consts.ScheduleCanCrawl.Code)
				if changed {
					logrus.Infoln("flow discover game", discoverGame.DoubanId)
				}
			}
		}
	}

	pendingSong := dao.SearchScheduleByStatus(consts.TypeSong.Code, consts.ScheduleToCrawl.Code)
	if pendingSong == nil {
		retrySong := dao.SearchScheduleByAll(consts.TypeSong.Code, consts.ScheduleCrawled.Code, consts.ScheduleUnready.Code)
		if retrySong != nil {
			changed := dao.CasScheduleStatus(retrySong.DoubanId, retrySong.Type, consts.ScheduleToCrawl.Code, consts.ScheduleCrawled.Code)
			if changed {
				logrus.Infoln("flow retry song ", retrySong.DoubanId)
			}
		} else {
			discoverSong := dao.SearchScheduleByStatus(consts.TypeSong.Code, consts.ScheduleCanCrawl.Code)
			if discoverSong != nil {
				changed := dao.CasScheduleStatus(discoverSong.DoubanId, discoverSong.Type, consts.ScheduleToCrawl.Code, consts.ScheduleCanCrawl.Code)
				if changed {
					logrus.Infoln("flow discover song", discoverSong.DoubanId)
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
				logrus.Infoln("flow retry user ", retryUser.DoubanId)
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
		for {
			runFlow()
		}
	}()

	logrus.Infoln("flow agent enabled")
}
