package agent

import (
	"mouban/consts"
	"mouban/crawl"
	"mouban/util"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func runLatest() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("run latest panic", r, "latest agent crashed  => ", util.GetCurrentGoroutineStack())
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
		for range time.NewTicker(time.Hour * 6).C {
			runLatest()
		}
	}()

	logrus.Infoln("latest agent enabled")
}
