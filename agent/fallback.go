package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mouban/consts"
	"mouban/dao"
	"mouban/util"
	"time"
)

func runFallback() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("run fallback panic", r, "fallback agent crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	cnt := dao.CasOrphanSchedule(consts.TypeUser.Code, time.Hour*6)
	logrus.Infoln(cnt, "orphan users reset")

	cnt = dao.CasOrphanSchedule(consts.TypeBook.Code, time.Hour*6)
	logrus.Infoln(cnt, "orphan books reset")

	cnt = dao.CasOrphanSchedule(consts.TypeMovie.Code, time.Hour*6)
	logrus.Infoln(cnt, "orphan movies reset")

	cnt = dao.CasOrphanSchedule(consts.TypeGame.Code, time.Hour*6)
	logrus.Infoln(cnt, "orphan games reset")

	cnt = dao.CasOrphanSchedule(consts.TypeSong.Code, time.Hour*6)
	logrus.Infoln(cnt, "orphan songs reset")

}

func init() {
	if !viper.GetBool("agent.enable") {
		logrus.Infoln("fallback agent disabled")
		return
	}

	go func() {
		for range time.NewTicker(time.Hour).C {
			runFallback()
		}
	}()

	logrus.Infoln("fallback agent enabled")
}
