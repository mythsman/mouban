package model

import (
	"time"
)

type Rating struct {
	ID        uint64
	Type      uint8  `gorm:"not null;uniqueIndex:uk_unique_id"`
	DoubanId  uint64 `gorm:"not null;uniqueIndex:uk_unique_id"`
	Total     uint64
	Rating    float32
	Star5     float32
	Star4     float32
	Star3     float32
	Star2     float32
	Star1     float32
	Status    uint8 `gorm:"comment:0-normal,1-not enough,2-can not rate;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
