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

// SearchScheduleByStatus idx_status
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

// SearchScheduleByAll idx_search
func SearchScheduleByAll(t uint8, status uint8, result uint8) *model.Schedule {
	schedule := &model.Schedule{}
	common.Db.Where("type = ? AND `status`= ? AND result = ?", t, status, result).
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

func CreateScheduleNx(doubanId uint64, t uint8, status uint8, result uint8) bool {
	existingSchedule := GetSchedule(doubanId, t)
	if existingSchedule == nil {
		schedule := &model.Schedule{
			DoubanId: doubanId,
			Type:     t,
			Status:   status,
			Result:   result,
		}
		row := common.Db.Create(&schedule).RowsAffected
		return row > 0
	} else {
		return false
	}
}
