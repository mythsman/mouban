package model

import (
	"time"
)

type Game struct {
	ID          uint64
	DoubanId    uint64 `gorm:"not null;uniqueIndex"`
	Title       string `gorm:"not null;type:varchar(512)"`
	Platform    string `gorm:"type:varchar(512)"`
	Genre       string `gorm:"type:varchar(512)"`
	Alias       string `gorm:"type:varchar(512)"`
	Developer   string `gorm:"type:varchar(512)"`
	Publisher   string `gorm:"type:varchar(512)"`
	PublishDate string `gorm:"type:varchar(512)"`
	Intro       string `gorm:"type:mediumtext"`
	Thumbnail   string `gorm:"type:varchar(512)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Game) TableName() string {
	return "game"
}

func (game Game) Show() *GameVO {
	return &GameVO{
		DoubanId:    game.DoubanId,
		Title:       game.Title,
		Platform:    game.Platform,
		Genre:       game.Genre,
		Alias:       game.Alias,
		Developer:   game.Developer,
		Publisher:   game.Publisher,
		PublishDate: game.PublishDate,
		Thumbnail:   game.Thumbnail,
	}
}

type GameVO struct {
	DoubanId    uint64 `json:"douban_id"`
	Title       string `json:"title"`
	Platform    string `json:"platform"`
	Genre       string `json:"genre"`
	Alias       string `json:"alias"`
	Developer   string `json:"developer"`
	Publisher   string `json:"publisher"`
	PublishDate string `json:"publish_date"`
	Thumbnail   string `json:"thumbnail"`
}
