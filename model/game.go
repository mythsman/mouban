package model

import (
	"time"
)

type Game struct {
	ID          uint64 `gorm:"primarykey"`
	DoubanId    uint64 `gorm:"uniqueIndex"`
	Title       string
	Platform    string
	Alias       string
	Developer   string
	Publisher   string
	PublishDate string
	Intro       string
	Thumbnail   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
