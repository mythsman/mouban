package model

import (
	"time"
)

type Movie struct {
	ID        uint64 `gorm:"primarykey"`
	DoubanId  uint64 `gorm:"uniqueIndex"`
	Title     string
	Director  string
	Writer    string
	Actor     string
	Style     string
	Site      string
	Country   string
	Language  string
	PublishAt string
	Season    uint64
	Episode   uint64
	Duration  string
	Alias     string
	IMDb      string
	Intro     string
	Thumbnail string
	CreatedAt time.Time
	UpdatedAt time.Time
}
