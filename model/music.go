package model

import (
	"time"
)

type Music struct {
	ID          uint64 `gorm:"primarykey"`
	DoubanId    uint64 `gorm:"uniqueIndex"`
	Title       string
	Actor       string
	Type        string
	Media       string
	Genre       string
	PublishDate string
	Publisher   string
	Barcode     string
	ISRC        string
	Alias       string
	Thumbnail   string
	Intro       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
