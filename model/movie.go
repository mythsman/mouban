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
	Actor       string `gorm:"type:varchar(2048)"`
	Style       string `gorm:"type:varchar(512)"`
	Site        string `gorm:"type:varchar(512)"`
	Country     string `gorm:"type:varchar(512)"`
	Language    string `gorm:"type:varchar(512)"`
	PublishDate string `gorm:"type:varchar(512)"`
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

func (movie Movie) Show() *MovieVO {
	return &MovieVO{
		DoubanId:    movie.DoubanId,
		Title:       movie.Title,
		Director:    movie.Director,
		Writer:      movie.Writer,
		Actor:       movie.Actor,
		PublishDate: movie.PublishDate,
		Alias:       movie.Alias,
		Thumbnail:   movie.Thumbnail,
	}
}

type MovieVO struct {
	DoubanId    uint64 `json:"douban_id"`
	Title       string `json:"title"`
	Director    string `json:"director"`
	Writer      string `json:"writer"`
	Actor       string `json:"actor"`
	PublishDate string `json:"publish_date"`
	Alias       string `json:"alias"`
	Thumbnail   string `json:"thumbnail"`
}
