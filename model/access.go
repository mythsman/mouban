package model

import (
	"time"
)

type Access struct {
	ID        uint64
	DoubanUid uint64 `gorm:"not null;index;"`
	Path      string `gorm:"type:varchar(64);"`
	Ip        string `gorm:"type:varchar(64);index"`
	UserAgent string `gorm:"type:varchar(512);"`
	Referer   string `gorm:"type:varchar(512);"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
