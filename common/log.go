package common

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)



func InitLogger() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableQuote:    true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.AddHook(&InfluxHooker{})

	influxdbUrl := viper.GetString("influxdb.url")
	if influxdbUrl == "" {
		return
	}

	client := influxdb2.NewClient(viper.GetString("influxdb.url"), viper.GetString("influxdb.token"))

	writeApi = client.WriteAPI(viper.GetString("influxdb.org"), viper.GetString("influxdb.bucket"))

	logrus.Infoln("influxdb init success")
}

type InfluxHooker struct {
}

var writeApi api.WriteAPI

func (h *InfluxHooker) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *InfluxHooker) Fire(e *logrus.Entry) error {
	if writeApi == nil {
		return nil
	}

	// Create point using fluent style
	p := influxdb2.NewPointWithMeasurement("log").
		AddTag("level", e.Level.String()).
		AddField("msg", e.Message).
		SetTime(time.Now())

	writeApi.WritePoint(p)

	return nil
}
