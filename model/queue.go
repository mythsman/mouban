package model

import (
	"time"
)

type Queue struct {
	ID        uint64
	DoubanUid uint64 `gorm:"not null;uniqueIndex"`
	Status    uint8  `gorm:"not null;"`//0-crawl done,1-need crawl,2-data invalid,3-unexpected fail
	CreatedAt time.Time
	UpdatedAt time.Time
}
