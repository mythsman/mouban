package model

import (
	"time"
)

type User struct {
	ID           uint64 `gorm:"primarykey"`
	DoubanUid    uint64
	UniqueId     string
	Name         string
	Thumbnail    string
	BookWish     uint64
	BookDo       uint64
	BookCollect  uint64
	GameWish     uint64
	GameDo       uint64
	GameCollect  uint64
	MusicWish    uint64
	MusicDo      uint64
	MusicCollect uint64
	MovieWish    uint64
	MovieDo      uint64
	MovieCollect uint64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
