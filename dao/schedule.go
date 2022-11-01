package dao

import (
	"fmt"
	"mouban/common"
	"mouban/consts"
	"mouban/model"
)

func GetSchedule(doubanId uint64, t uint8) *model.Schedule {
	schedule := &model.Schedule{}
	common.Db.Where("douban_id = ? AND type = ? ", doubanId, t).Find(schedule)
	return schedule
}

func SearchSchedule(t uint8, status uint8, limit int) *[]model.Schedule {
	var schedules []model.Schedule
	common.Db.Limit(limit).Where("type = ? AND status = ? ", t, status).Find(&schedules)
	return &schedules
}

func UpsertSchedule(doubanId uint64, t uint8, status uint8) {
	if t != consts.TypeBook &&
		t != consts.TypeMovie &&
		t != consts.TypeGame &&
		t != consts.TypeUser {
		fmt.Println("type invalid : ", t)
		return
	}
	if status != consts.ScheduleToCrawl &&
		status != consts.ScheduleCrawling &&
		status != consts.ScheduleSucceeded &&
		status != consts.ScheduleFailed &&
		status != consts.ScheduleInvalid {
		fmt.Println("status invalid : ", status)
	}

	schedule := &model.Schedule{
		DoubanId: doubanId,
		Type:     t,
		Status:   status,
	}
	if common.Db.Where("douban_id = ? AND type = ? ", doubanId, t).Updates(&schedule).RowsAffected == 0 {
		common.Db.Create(&schedule)
	}
}
