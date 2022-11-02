package dao

import (
	"fmt"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func TestUpsertSchedule(t *testing.T) {
	UpsertSchedule(1234, consts.TypeMovie, consts.ScheduleStatusToCrawl, consts.ScheduleResultInvalid)
}

func TestGetSchedule(t *testing.T) {
	schedule := GetSchedule(1234, consts.TypeMovie)
	fmt.Println(util.ToJson(schedule))
}

func TestSearchSchedule(t *testing.T) {
	schedules := SearchScheduleByStatus(consts.TypeMovie, consts.ScheduleStatusToCrawl, 10)
	fmt.Println(util.ToJson(schedules))
}
