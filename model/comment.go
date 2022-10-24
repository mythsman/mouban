package model

import (
	"time"
)

type Comment struct {
	ID        uint64 `gorm:"primarykey"`
	DoubanId  uint64
	Type      uint8
	Rate      uint8
	DoubanUid uint64
	Label     string
	Comment   string
	Status    uint8
	MarkDate  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
