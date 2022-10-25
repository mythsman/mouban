package model

import (
	"time"
)

type User struct {
	ID           uint64
	DoubanUid    uint64 `gorm:"not null;uniqueIndex"`
	UniqueId     string `gorm:"not hull;type:varchar(64);uniqueIndex"`
	Name         string `gorm:"not null;type:varchar(512);"`
	Thumbnail    string `gorm:"type:varchar(512);"`
	BookWish     uint32
	BookDo       uint32
	BookCollect  uint32
	GameWish     uint32
	GameDo       uint32
	GameCollect  uint32
	MusicWish    uint32
	MusicDo      uint32
	MusicCollect uint32
	MovieWish    uint32
	MovieDo      uint32
	MovieCollect uint32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (User) TableName() string {
	return "user"
}
