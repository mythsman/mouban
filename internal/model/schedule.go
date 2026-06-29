package model

import (
	"time"
)

type Schedule struct {
	ID        uint64
	DoubanId  uint64 `gorm:"not null;uniqueIndex:uk_schedule;"`
	Type      uint8  `gorm:"not null;uniqueIndex:uk_schedule;index:idx_status,priority:1;index:idx_result,priority:1;index:idx_search,priority:1"`
	Status    *uint8 `gorm:"not null;index:idx_status,priority:2;index:idx_search,priority:2;index:idx_schedule_status_updated,priority:1"`
	Result    *uint8 `gorm:"not null;index:idx_result,priority:2;index:idx_search,priority:2"`
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"index:idx_result,priority:3;index:idx_status,priority:3;index:idx_search,priority:3;index:idx_schedule_status_updated,priority:2"`
}

func (Schedule) TableName() string {
	return "schedule"
}
