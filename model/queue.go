package model

import (
	"time"
)

type Queue struct {
	ID        uint64 `gorm:"primarykey"`
	DoubanUid uint64 `gorm:"uniqueIndex"`
	Status    uint8  //0-crawl done,1-need crawl,2-data invalid
	CreatedAt time.Time
	UpdatedAt time.Time
}
