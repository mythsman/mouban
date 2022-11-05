package dao

import (
	"log"
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

func SearchScheduleByStatus(status uint8) *model.Schedule {
	schedule := &model.Schedule{}
	common.Db.Where("status = ? ", status).
		Order("updated_at asc").
		Find(&schedule)
	if schedule.ID == 0 {
		return nil
	}
	return schedule
}

func SearchScheduleByResult(t uint8, result uint8, limit int) *[]model.Schedule {
	var schedules []model.Schedule
	common.Db.Limit(limit).Where("type = ? AND result = ? ", t, result).Find(&schedules)
	return &schedules
}

func CasScheduleStatus(doubanId uint64, t uint8, status uint8, rawStatus uint8) bool {
	row := common.Db.Model(&model.Schedule{}).
		Where("douban_id = ? AND type = ? AND status = ?", doubanId, t, rawStatus).
		Update("status", status).RowsAffected
	return row > 0
}

func ChangeScheduleResult(doubanId uint64, t uint8, result uint8) {
	common.Db.Model(&model.Schedule{}).
		Where("douban_id = ? AND type = ?", doubanId, t).
		Update("result", result)
}

func CreateSchedule(doubanId uint64, t uint8, status uint8, result uint8) {
	if t != consts.TypeBook &&
		t != consts.TypeMovie &&
		t != consts.TypeGame &&
		t != consts.TypeUser {
		log.Println("type invalid : ", t)
		return
	}
	if status != consts.ScheduleStatusToCrawl &&
		status != consts.ScheduleStatusCrawling &&
		status != consts.ScheduleStatusCrawled {
		log.Println("status invalid : ", status)
		return
	}

	if result != consts.ScheduleResultUnready &&
		result != consts.ScheduleResultReady &&
		result != consts.ScheduleResultInvalid {
		log.Println("result invalid : ", status)
		return
	}

	schedule := &model.Schedule{
		DoubanId: doubanId,
		Type:     t,
		Status:   status,
		Result:   result,
	}
	common.Db.Create(&schedule)
}
