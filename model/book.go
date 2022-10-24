package model

import (
	"time"
)

type Book struct {
	ID          uint64 `gorm:"primarykey"`
	DoubanId    uint64 `gorm:"uniqueIndex"`
	Title       string
	Subtitle    string
	Author      string
	Translator  string
	Publisher   string
	PublishDate string
	ISBN        string
	Page        string
	Price       uint64
	Intro       string
	Thumbnail   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
