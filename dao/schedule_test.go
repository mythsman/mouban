package dao

import (
	"fmt"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func TestUpsertSchedule(t *testing.T) {
	CreateSchedule(1234, consts.TypeMovie, consts.ScheduleStatusToCrawl, consts.ScheduleResultInvalid)
}

func TestGetSchedule(t *testing.T) {
	schedule := GetSchedule(1234, consts.TypeMovie)
	fmt.Println(util.ToJson(schedule))
}

func TestSearchSchedule(t *testing.T) {
	schedule := SearchScheduleByStatus(consts.TypeBook, consts.ScheduleStatusToCrawl)
	fmt.Println(util.ToJson(schedule))
}

func TestChangeScheduleResult(t *testing.T) {
	ChangeScheduleResult(162448367, consts.TypeUser, consts.ScheduleResultReady)
}
