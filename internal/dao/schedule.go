package dao

import (
	"mouban/internal/common"
	"mouban/internal/consts"
	"mouban/internal/model"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	dataProcessTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mouban_data_process_total",
		Help: "Data processed counter",
	}, []string{"type", "result"})
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
	dataProcessTotal.WithLabelValues(consts.ParseType(t).Name, consts.ParseResult(result).Name).Inc()

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

func ListScheduleByStatus(status uint8, limit int) []model.Schedule {
	if limit <= 0 {
		limit = 50
	}
	rows := make([]model.Schedule, 0)
	common.Db.Where("status = ?", status).
		Order("updated_at asc").
		Limit(limit).
		Find(&rows)
	return rows
}

func ListRecentScheduleByStatus(status uint8, limit int) []model.Schedule {
	if limit <= 0 {
		limit = 50
	}
	rows := make([]model.Schedule, 0)
	common.Db.Where("status = ?", status).
		Order("updated_at desc").
		Limit(limit).
		Find(&rows)
	return rows
}

type typeStatusCountRow struct {
	Type   uint8 `gorm:"column:type"`
	Status uint8 `gorm:"column:status"`
	Count  int64 `gorm:"column:count"`
}

func CountScheduleByStatusesGrouped(statuses []uint8) map[uint8]map[uint8]int64 {
	result := map[uint8]map[uint8]int64{}
	if len(statuses) == 0 {
		return result
	}

	// []uint8 is alias of []byte; for SQL IN we must avoid passing binary payload.
	statusList := make([]int, 0, len(statuses))
	for _, s := range statuses {
		statusList = append(statusList, int(s))
	}

	rows := make([]typeStatusCountRow, 0)
	common.Db.Model(&model.Schedule{}).
		Select("type, status, COUNT(*) as count").
		Where("status IN ?", statusList).
		Group("type, status").
		Find(&rows)

	for _, row := range rows {
		if result[row.Type] == nil {
			result[row.Type] = map[uint8]int64{}
		}
		result[row.Type][row.Status] = row.Count
	}
	return result
}

type typeOldestRow struct {
	Type   uint8     `gorm:"column:type"`
	Oldest time.Time `gorm:"column:oldest"`
}

func FindOldestUpdatedAtByStatusGrouped(status uint8) map[uint8]time.Time {
	result := map[uint8]time.Time{}
	rows := make([]typeOldestRow, 0)
	common.Db.Model(&model.Schedule{}).
		Select("type, MIN(updated_at) as oldest").
		Where("status = ?", status).
		Group("type").
		Scan(&rows)

	for _, row := range rows {
		result[row.Type] = row.Oldest
	}
	return result
}
