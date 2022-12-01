package dao

import (
	"mouban/common"
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

	schedule := &model.Schedule{
		DoubanId: doubanId,
		Type:     t,
		Status:   status,
		Result:   result,
	}
	common.Db.Create(&schedule)
}
