package model

import (
	"time"
)

type Rating struct {
	ID        uint64 `gorm:"primarykey"`
	Type      uint8
	DoubanId  uint64
	Total     uint64
	Rating    float32
	Star5     float32
	Star4     float32
	Star3     float32
	Star2     float32
	Star1     float32
	CreatedAt time.Time
	UpdatedAt time.Time
}
