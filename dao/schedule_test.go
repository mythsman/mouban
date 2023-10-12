package dao

import (
	"mouban/consts"
	"mouban/util"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestUpsertSchedule(t *testing.T) {
	CreateScheduleNx(1234, consts.TypeMovie.Code, consts.ScheduleToCrawl.Code, consts.ScheduleInvalid.Code)
}

func TestCasOrphanSchedule(t *testing.T) {
	cnt := CasOrphanSchedule(consts.TypeBook.Code, 0)
	logrus.Infoln(cnt)
}

func TestGetSchedule(t *testing.T) {
	schedule := GetSchedule(1234, consts.TypeMovie.Code)
	logrus.Infoln(util.ToJson(schedule))
}

func TestSearchSchedule(t *testing.T) {
	schedule := SearchScheduleByStatus(consts.TypeBook.Code, consts.ScheduleToCrawl.Code)
	logrus.Infoln(util.ToJson(schedule))
}

func TestChangeScheduleResult(t *testing.T) {
	ChangeScheduleResult(162448367, consts.TypeUser.Code, consts.ScheduleReady.Code)
}
