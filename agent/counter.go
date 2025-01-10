package agent

import (
	"mouban/consts"
	"mouban/dao"
	"mouban/util"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

var (
	dataCountGuage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "mouban_data_count",
		Help: "Diff data count",
	}, []string{"type"})
)

func runCounter() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("run counter panic", r, "counter agent crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()

	dataCountGuage.WithLabelValues(consts.TypeBook.Name).Set(float64(dao.CountBook()))
	dataCountGuage.WithLabelValues(consts.TypeMovie.Name).Set(float64(dao.CountMovie()))
	dataCountGuage.WithLabelValues(consts.TypeGame.Name).Set(float64(dao.CountGame()))
	dataCountGuage.WithLabelValues(consts.TypeSong.Name).Set(float64(dao.CountSong()))
	dataCountGuage.WithLabelValues(consts.TypeUser.Name).Set(float64(dao.CountUser()))
}

func init() {
	runCounter()

	go func() {
		for range time.NewTicker(time.Hour).C {
			runCounter()
		}
	}()

	logrus.Infoln("counter agent enabled")
}
