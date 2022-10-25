package model

import (
	"time"
)

type Comment struct {
	ID        uint64
	DoubanUid uint64 `gorm:"not null;uniqueIndex:uk_comment;index:idx_comment;priority:1"`
	DoubanId  uint64 `gorm:"not null;uniqueIndex:uk_comment;index:idx_comment;priority:2"`
	Type      uint8  `gorm:"not null;uniqueIndex:uk_comment;index:idx_comment;priority:3"`
	Rate      uint8
	Label     string    `gorm:"type:varchar(512)"`
	Comment   string    `gorm:"type:mediumtext"`
	Status    uint8     `gorm:"not null;index:idx_comment;priority:4"`
	MarkDate  time.Time `gorm:"not null;index:idx_comment;priority:5"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
