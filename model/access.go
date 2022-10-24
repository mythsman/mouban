package model

import (
	"time"
)

type Access struct {
	ID        uint64
	DoubanUid uint64 `gorm:"index:"`
	Path      string
	Ip        string `gorm:"index:"`
	UserAgent string
	Referer   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
