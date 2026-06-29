package common

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// gormStructuredLogger outputs slow-query logs in structured fields.
type gormStructuredLogger struct {
	slowThreshold             time.Duration
	logLevel                  logger.LogLevel
	ignoreRecordNotFoundError bool
}

func newGormStructuredLogger(slowThreshold time.Duration) logger.Interface {
	return &gormStructuredLogger{
		slowThreshold:             slowThreshold,
		logLevel:                  logger.Warn,
		ignoreRecordNotFoundError: true,
	}
}

func (l *gormStructuredLogger) LogMode(level logger.LogLevel) logger.Interface {
	clone := *l
	clone.logLevel = level
	return &clone
}

func (l *gormStructuredLogger) Info(_ context.Context, msg string, data ...interface{}) {
	if l.logLevel < logger.Info {
		return
	}
	logrus.WithFields(logrus.Fields{
		"component": "gorm",
		"message":   msg,
		"data":      data,
	}).Info("gorm info")
}

func (l *gormStructuredLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	if l.logLevel < logger.Warn {
		return
	}
	logrus.WithFields(logrus.Fields{
		"component": "gorm",
		"message":   msg,
		"data":      data,
	}).Warn("gorm warn")
}

func (l *gormStructuredLogger) Error(_ context.Context, msg string, data ...interface{}) {
	if l.logLevel < logger.Error {
		return
	}
	logrus.WithFields(logrus.Fields{
		"component": "gorm",
		"message":   msg,
		"data":      data,
	}).Error("gorm error")
}

func (l *gormStructuredLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel == logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil && l.logLevel >= logger.Error {
		if !(errors.Is(err, gorm.ErrRecordNotFound) && l.ignoreRecordNotFoundError) {
			logrus.WithFields(logrus.Fields{
				"component":  "mysql",
				"elapsed_ms": elapsed.Milliseconds(),
				"rows":       rows,
				"sql":        sql,
				"error":      err.Error(),
			}).Error("mysql query failed")
		}
		return
	}

	if elapsed > l.slowThreshold && l.logLevel >= logger.Warn {
		logrus.WithFields(logrus.Fields{
			"component":    "mysql",
			"elapsed_ms":   elapsed.Milliseconds(),
			"threshold_ms": l.slowThreshold.Milliseconds(),
			"rows":         rows,
			"sql":          sql,
		}).Warn("mysql slow query")
	}
}
