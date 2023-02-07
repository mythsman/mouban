package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var logger *logrus.Logger

func Info(args ...interface{}) {
	logger.Info(args)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

func Instance() *logrus.Logger {
	return logger
}

func init() {
	logger = logrus.New()
	logger.Out = os.Stdout
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableQuote:    true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
