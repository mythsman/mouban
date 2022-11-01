package dao

import (
	"fmt"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func TestChangeSchedule(t *testing.T) {
	ChangeSchedule(1234, consts.TypeMovie, consts.ScheduleSucceeded)
}

func TestGetSchedule(t *testing.T) {
	schedule := GetSchedule(1234, consts.TypeMovie)
	fmt.Println(util.ToJson(schedule))
}

func TestSearchSchedule(t *testing.T) {
	schedules := SearchSchedule(consts.TypeMovie, consts.ScheduleSucceeded, 10)
	fmt.Println(util.ToJson(schedules))
}
