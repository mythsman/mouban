package model

import (
	"time"
)

type Record struct {
	ID        uint64
	DoubanId  uint64 `gorm:"not null;uniqueIndex:uk_record;"`
	Type      uint8  `gorm:"not null;uniqueIndex:uk_record;index:idx_search;priority=1"`
	Status    uint8  `gorm:"not null;index:idx_search;priority=2"`
	TimeCost  uint32
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"index:idx_search;priority=3"`
}

func (Record) TableName() string {
	return "record"
}
