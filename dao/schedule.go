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
	if schedule.ID == 0 {
		return nil
	}
	return schedule
}

func SearchScheduleByStatus(t uint8, status uint8, limit int) *[]model.Schedule {
	var schedules []model.Schedule
	common.Db.Limit(limit).Where("type = ? AND status = ? ", t, status).Find(&schedules)
	return &schedules
}

func SearchScheduleByResult(t uint8, result uint8, limit int) *[]model.Schedule {
	var schedules []model.Schedule
	common.Db.Limit(limit).Where("type = ? AND result = ? ", t, result).Find(&schedules)
	return &schedules
}

func UpsertSchedule(doubanId uint64, t uint8, status uint8, result uint8) {
	if t != consts.TypeBook &&
		t != consts.TypeMovie &&
		t != consts.TypeGame &&
		t != consts.TypeUser {
		fmt.Println("type invalid : ", t)
		return
	}
	if status != consts.ScheduleStatusToCrawl &&
		status != consts.ScheduleStatusCrawling &&
		status != consts.ScheduleStatusCrawled {
		fmt.Println("status invalid : ", status)
		return
	}

	if result != consts.ScheduleResultUnready &&
		result != consts.ScheduleResultReady &&
		result != consts.ScheduleResultInvalid {
		fmt.Println("result invalid : ", status)
		return
	}

	schedule := &model.Schedule{
		DoubanId: doubanId,
		Type:     t,
		Status:   status,
		Result:   result,
	}
	if common.Db.Where("douban_id = ? AND type = ? ", doubanId, t).Updates(&schedule).RowsAffected == 0 {
		common.Db.Create(&schedule)
	}
}
