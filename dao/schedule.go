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

func SearchScheduleByStatus(t uint8, status uint8) *model.Schedule {
	schedule := &model.Schedule{}
	common.Db.Where("type = ? AND status = ? ", t, status).
		Order("updated_at asc").
		Limit(1).
		Find(&schedule)
	if schedule.ID == 0 {
		return nil
	}
	return schedule
}

func SearchSchedule(status uint8, result uint8) *model.Schedule {
	schedule := &model.Schedule{}
	common.Db.Where("`status`= ? AND result = ?", status, result).
		Order("updated_at asc").
		Limit(1).
		Find(&schedule)
	if schedule.ID == 0 {
		return nil
	}
	return schedule
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
