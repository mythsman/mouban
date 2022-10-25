package model

import (
	"time"
)

type Movie struct {
	ID        uint64
	DoubanId  uint64 `gorm:"not null;uniqueIndex"`
	Title     string `gorm:"not null;type:varchar(512)"`
	Director  string `gorm:"type:varchar(512)"`
	Writer    string `gorm:"type:varchar(512)"`
	Actor     string `gorm:"type:varchar(512)"`
	Style     string `gorm:"type:varchar(512)"`
	Site      string `gorm:"type:varchar(512)"`
	Country   string `gorm:"type:varchar(512)"`
	Language  string `gorm:"type:varchar(512)"`
	PublishAt string `gorm:"type:datetime(3)"`
	Season    uint32
	Episode   uint32
	Duration  uint32
	Alias     string `gorm:"type:varchar(512)"`
	IMDb      string `gorm:"type:varchar(512)"`
	Intro     string `gorm:"type:mediumtext"`
	Thumbnail string `gorm:"type:varchar(512)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Movie) TableName() string {
	return "movie"
}
