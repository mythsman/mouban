package model

import (
	"time"
)

type Book struct {
	ID          uint64
	DoubanId    uint64 `gorm:"not null;uniqueIndex"`
	Title       string `gorm:"not null;type:varchar(1024)"`
	Subtitle    string `gorm:"type:varchar(1024)"`
	Orititle    string `gorm:"type:varchar(1024)"`
	Author      string `gorm:"type:varchar(1024)"`
	Translator  string `gorm:"type:varchar(512)"`
	Press       string `gorm:"type:varchar(512)"`
	Producer    string `gorm:"type:varchar(512)"`
	Serial      string `gorm:"type:varchar(512)"`
	PublishDate string `gorm:"type:varchar(64)"`
	ISBN        string `gorm:"type:varchar(64)"`
	Framing     string `gorm:"type:varchar(512)"`
	Page        uint32
	Price       uint32
	BookIntro   string `gorm:"type:mediumtext"`
	AuthorIntro string `gorm:"type:mediumtext"`
	Thumbnail   string `gorm:"type:varchar(512)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Book) TableName() string {
	return "book"
}

func (book Book) Show() *BookVO {
	return &BookVO{
		DoubanId:    book.DoubanId,
		Title:       book.Title,
		Subtitle:    book.Subtitle,
		Orititle:    book.Orititle,
		Author:      book.Author,
		Translator:  book.Translator,
		Press:       book.Press,
		Producer:    book.Producer,
		PublishDate: book.PublishDate,
		Thumbnail:   book.Thumbnail,
	}
}

type BookVO struct {
	DoubanId    uint64 `json:"douban_id"`
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Orititle    string `json:"orititle"`
	Author      string `json:"author"`
	Translator  string `json:"translator"`
	Press       string `json:"press"`
	Producer    string `json:"producer"`
	PublishDate string `json:"publish_date"`
	Thumbnail   string `json:"thumbnail"`
}
