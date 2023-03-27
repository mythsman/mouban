package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mouban/consts"
	"mouban/crawl"
	"mouban/util"
	"time"
)

func runLatest() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, "latest agent crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	bookIds := crawl.LatestBook()
	logrus.Infoln(len(*bookIds), "latest books discovered")
	processDiscoverItem(bookIds, consts.TypeBook)

	movieIds := crawl.LatestMovie()
	logrus.Infoln(len(*movieIds), "latest movies discovered")
	processDiscoverItem(movieIds, consts.TypeMovie)

	songIds := crawl.LatestSong()
	logrus.Infoln(len(*songIds), "latest songs discovered")
	processDiscoverItem(songIds, consts.TypeSong)

}

func init() {
	if !viper.GetBool("agent.enable") {
		logrus.Infoln("latest agent disabled")
		return
	}

	go func() {
		for range time.NewTicker(time.Hour * 24).C {
			runLatest()
		}
	}()

	logrus.Infoln("latest agent enabled")
}
