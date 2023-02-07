package dao

import (
	"github.com/sirupsen/logrus"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func TestUpsertSchedule(t *testing.T) {
	CreateScheduleNx(1234, consts.TypeMovie.Code, consts.ScheduleStatusToCrawl, consts.ScheduleResultInvalid)
}

func TestCasOrphanSchedule(t *testing.T) {
	cnt := CasOrphanSchedule(0)
	logrus.Info(cnt)
}

func TestGetSchedule(t *testing.T) {
	schedule := GetSchedule(1234, consts.TypeMovie.Code)
	logrus.Info(util.ToJson(schedule))
}

func TestSearchSchedule(t *testing.T) {
	schedule := SearchScheduleByStatus(consts.TypeBook.Code, consts.ScheduleStatusToCrawl)
	logrus.Info(util.ToJson(schedule))
}

func TestChangeScheduleResult(t *testing.T) {
	ChangeScheduleResult(162448367, consts.TypeUser.Code, consts.ScheduleResultReady)
}
