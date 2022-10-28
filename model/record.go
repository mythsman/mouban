package model

import (
	"time"
)

type Record struct {
	ID        uint64
	DoubanId  uint64 `gorm:"not null;uniqueIndex:uk_record;"`
	Type      uint8  `gorm:"not null;uniqueIndex:uk_record;index:idx_search;priority=1"`
	Status    uint8  `gorm:"not null;index:idx_search;priority=2"` //0-crawl succeeded,1-crawling,2-crawl failed,3-data invalid
	CreatedAt time.Time
	UpdatedAt time.Time `gorm:"index:idx_search;priority=3"`
}

func (Record) TableName() string {
	return "record"
}
