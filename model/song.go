package model

import (
	"time"
)

type Song struct {
	ID          uint64
	DoubanId    uint64 `gorm:"not null;uniqueIndex"`
	Title       string `gorm:"not null;type:varchar(512)"`
	Alias       string `gorm:"type:varchar(512)"`
	Musician    string `gorm:"type:varchar(2048)"`
	AlbumType   string `gorm:"type:varchar(512)"` //专辑类型
	Genre       string `gorm:"type:varchar(512)"` //流派
	Media       string `gorm:"type:varchar(512)"`
	Barcode     string `gorm:"type:varchar(512)"` //条形码
	Publisher   string `gorm:"type:varchar(512)"`
	PublishDate string `gorm:"type:varchar(512)"`
	ISRC        string `gorm:"type:varchar(512)"`
	AlbumCount  uint32 //唱片数量
	Intro       string `gorm:"type:mediumtext"`
	TrackList   string `gorm:"type:mediumtext"`
	Thumbnail   string `gorm:"type:varchar(512)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Song) TableName() string {
	return "song"
}

func (song Song) Show() *SongVO {
	return &SongVO{
		DoubanId:    song.DoubanId,
		Title:       song.Title,
		Alias:       song.Alias,
		Musician:    song.Musician,
		Thumbnail:   song.Thumbnail,
		AlbumType:   song.AlbumType,
		Genre:       song.Genre,
		Media:       song.Media,
		Publisher:   song.Publisher,
		PublishDate: song.PublishDate,
	}
}

type SongVO struct {
	DoubanId    uint64 `json:"douban_id"`
	Title       string `json:"title"`
	Alias       string `json:"alias"`
	Musician    string `json:"musician"`
	Thumbnail   string `json:"thumbnail"`
	AlbumType   string `json:"album_type"`
	Genre       string `json:"genre"`
	Media       string `json:"media"`
	Publisher   string `json:"publisher"`
	PublishDate string `json:"publish_date"`
}
