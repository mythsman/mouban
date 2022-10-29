package model

import (
	"time"
)

type Book struct {
	ID           uint64
	DoubanId     uint64 `gorm:"not null;uniqueIndex"`
	Title        string `gorm:"not null;type:varchar(512)"`
	Subtitle     string `gorm:"type:varchar(512)"`
	Author       string `gorm:"type:varchar(512)"`
	Translator   string `gorm:"type:varchar(512)"`
	Press        string `gorm:"type:varchar(512)"`
	Producer     string `gorm:"type:varchar(512)"`
	Serial       string `gorm:"type:varchar(512)"`
	PublishMonth time.Time
	ISBN         string `gorm:"type:varchar(64)"`
	Page         uint32
	Price        uint32
	Intro        string `gorm:"type:mediumtext"`
	Thumbnail    string `gorm:"type:varchar(512)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Book) TableName() string {
	return "book"
}

type BookVO struct {
	DoubanId     uint64   `json:"douban_id"`
	Title        string   `json:"title"`
	Subtitle     string   `json:"subtitle"`
	Author       []string `json:"author"`
	Translator   []string `json:"translator"`
	Press        string   `json:"press"`
	Producer     string   `json:"producer"`
	Serial       string   `json:"serial"`
	PublishMonth string   `json:"publish_month"`
	ISBN         string   `json:"isbn"`
	Page         uint32   `json:"page"`
	Price        uint32   `json:"price"`
	Intro        string   `json:"intro"`
	Thumbnail    string   `json:"thumbnail"`
}
