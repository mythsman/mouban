package model

import (
	"time"
)

type Movie struct {
	ID          uint64
	DoubanId    uint64 `gorm:"not null;uniqueIndex"`
	Title       string `gorm:"not null;type:varchar(512)"`
	Director    string `gorm:"type:varchar(512)"`
	Writer      string `gorm:"type:varchar(512)"`
	Actor       string `gorm:"type:varchar(1024)"`
	Style       string `gorm:"type:varchar(512)"`
	Site        string `gorm:"type:varchar(512)"`
	Country     string `gorm:"type:varchar(512)"`
	Language    string `gorm:"type:varchar(512)"`
	PublishDate string `gorm:"type:varchar(64)"`
	Episode     uint32
	Duration    uint32
	Alias       string `gorm:"type:varchar(512)"`
	IMDb        string `gorm:"type:varchar(512);column:imdb"`
	Intro       string `gorm:"type:mediumtext"`
	Thumbnail   string `gorm:"type:varchar(512)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Movie) TableName() string {
	return "movie"
}

type MovieVO struct {
	DoubanId    uint64   `json:"douban_id"`
	Title       string   `json:"title"`
	Director    []string `json:"director"`
	Writer      []string `json:"writer"`
	Actor       []string `json:"actor"`
	Style       []string `json:"style"`
	Site        string   `json:"site"`
	Country     []string `json:"country"`
	Language    []string `json:"language"`
	PublishDate []string `json:"publish_date"`
	Episode     uint32   `json:"episode"`
	Duration    uint32   `json:"duration"`
	Alias       string   `json:"alias"`
	IMDb        string   `json:"imdb"`
	Intro       string   `json:"intro"`
	Thumbnail   string   `json:"thumbnail"`
}
