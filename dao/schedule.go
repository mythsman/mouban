package dao

import (
	"mouban/common"
	"mouban/consts"
	"mouban/model"
	"time"
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

// CasOrphanSchedule idx_status
func CasOrphanSchedule(t uint8, expire time.Duration) int64 {
	return common.Db.Model(&model.Schedule{}).
		Where("type = ? AND status = ? AND updated_at < ?", t, consts.ScheduleCrawling.Code, time.Now().Add(-expire)).
		Update("status", consts.ScheduleToCrawl.Code).RowsAffected
}

// CasScheduleStatus uk_schedule
func CasScheduleStatus(doubanId uint64, t uint8, status uint8, rawStatus uint8) bool {
	row := common.Db.Model(&model.Schedule{}).
		Where("douban_id = ? AND type = ? AND status = ?", doubanId, t, rawStatus).
		Update("status", status).RowsAffected
	return row > 0
}

// ChangeScheduleResult uk_schedule
func ChangeScheduleResult(doubanId uint64, t uint8, result uint8) {
	common.Db.Model(&model.Schedule{}).
		Where("douban_id = ? AND type = ?", doubanId, t).
		Update("result", result)
}

func CreateScheduleNx(doubanId uint64, t uint8, status uint8, result uint8) bool {
	data := &model.Schedule{}
	insert := &model.Schedule{
		DoubanId: doubanId,
		Type:     t,
		Status:   &status,
		Result:   &result,
	}
	row := common.Db.Where("douban_id = ? AND type = ? ", doubanId, t).Attrs(insert).FirstOrCreate(data).RowsAffected
	return row > 0
}
