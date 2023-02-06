package log

import "github.com/sirupsen/logrus"

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
	logger.SetFormatter(&logrus.TextFormatter{
		DisableQuote: true,
	})
}
