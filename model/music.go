package model

import (
	"time"
)

type Music struct {
	ID          uint64
	DoubanId    uint64 `gorm:"not null;uniqueIndex"`
	Title       string `gorm:"not null;type:varchar(512)"`
	Actor       string `gorm:"type:varchar(512)"`
	Type        string `gorm:"type:varchar(512)"`
	Media       string `gorm:"type:varchar(512)"`
	Genre       string `gorm:"type:varchar(512)"`
	PublishDate string `gorm:"type:varchar(512)"`
	Publisher   string `gorm:"type:varchar(512)"`
	Barcode     string `gorm:"type:varchar(64)"`
	ISRC        string `gorm:"type:varchar(64)"`
	Alias       string `gorm:"type:varchar(512)"`
	Thumbnail   string `gorm:"type:varchar(512)"`
	Intro       string `gorm:"type:mediumtext"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
