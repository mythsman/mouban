package model

import (
	"time"
)

type Access struct {
	ID        uint64 `gorm:"primarykey"`
	DoubanUid uint64
	Path      string
	Ip        string
	UserAgent string
	Referer   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
