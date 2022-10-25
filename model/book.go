package model

import (
	"time"
)

type Book struct {
	ID          uint64
	DoubanId    uint64 `gorm:"not null;uniqueIndex"`
	Title       string `gorm:"type:varchar(512)"`
	Subtitle    string `gorm:"type:varchar(512)"`
	Author      string `gorm:"type:varchar(512)"`
	Translator  string `gorm:"type:varchar(512)"`
	Publisher   string `gorm:"type:varchar(512)"`
	PublishDate string `gorm:"type:varchar(512)"`
	ISBN        string `gorm:"type:varchar(64)"`
	Page        string `gorm:"type:varchar(64)"`
	Price       uint64
	Intro       string `gorm:"type:mediumtext"`
	Thumbnail   string `gorm:"type:varchar(512)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
