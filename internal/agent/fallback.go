package agent

import (
	"mouban/internal/consts"
	"mouban/internal/dao"
	"mouban/internal/util"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func runFallback() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("run fallback panic", r, "fallback agent crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	expire, err := time.ParseDuration(viper.GetString("agent.orphan_expire"))
	if err != nil || expire <= 0 {
		expire = time.Hour * 6
		logrus.Warnln("invalid agent.orphan_expire, fallback to", expire)
	}

	cnt := dao.CasOrphanSchedule(consts.TypeUser.Code, expire)
	logrus.Infoln(cnt, "orphan users reset")

	cnt = dao.CasOrphanSchedule(consts.TypeBook.Code, expire)
	logrus.Infoln(cnt, "orphan books reset")

	cnt = dao.CasOrphanSchedule(consts.TypeMovie.Code, expire)
	logrus.Infoln(cnt, "orphan movies reset")

	cnt = dao.CasOrphanSchedule(consts.TypeGame.Code, expire)
	logrus.Infoln(cnt, "orphan games reset")

	cnt = dao.CasOrphanSchedule(consts.TypeSong.Code, expire)
	logrus.Infoln(cnt, "orphan songs reset")
}

func startFallbackAgent() {
	if !viper.GetBool("agent.enable") {
		logrus.Infoln("fallback agent disabled")
		return
	}

	runFallback()

	go func() {
		for range time.NewTicker(time.Hour).C {
			runFallback()
		}
	}()

	logrus.Infoln("fallback agent enabled")
}
