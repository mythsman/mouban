package model

import (
	"time"
)

type Music struct {
	ID          uint64
	DoubanId    uint64 `gorm:"not null;uniqueIndex"`
	Title       string `gorm:"not null;type:varchar(512)"`
	Actor       string `gorm:"type:varchar(512)"`
	Style        string `gorm:"type:varchar(512)"`
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

func (Music) TableName() string {
	return "music"
}


type MusicVO struct {
	DoubanId    uint64 `json:"douban_id"`
	Title       string `json:"title"`
	Actor       []string `json:"actor"`
	Style       []string `json:"style"`
	Media       []string `json:"media"`
	Genre       []string `json:"genre"`
	PublishDate string `json:"publish_date"`
	Publisher   []string `json:"publisher"`
	Barcode     string `json:"barcode"`
	ISRC        string `json:"isrc"`
	Alias       string `json:"alias"`
	Thumbnail   string `json:"thumbnail"`
	Intro       string `json:"intro"`
}
