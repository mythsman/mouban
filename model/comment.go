package model

import (
	"time"
)

type Comment struct {
	ID        uint64 `gorm:"primarykey"`
	DoubanUid uint64 `gorm:"uniqueIndex:uk_comment,index:idx_comment,priority:1"`
	DoubanId  uint64 `gorm:"uniqueIndex:uk_comment,index:idx_comment,priority:2"`
	Type      uint8  `gorm:"uniqueIndex:uk_comment,index:idx_comment,priority:3"`
	Rate      uint8
	Label     string
	Comment   string
	Status    uint8     `gorm:"index:idx_comment,priority:4"`
	MarkDate  time.Time `gorm:"index:idx_comment,priority:5"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
