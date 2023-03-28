package model

import (
	"time"
)

type User struct {
	ID           uint64
	DoubanUid    uint64 `gorm:"not null;uniqueIndex"`
	Domain       string `gorm:"not hull;index;type:varchar(64)"`
	Name         string `gorm:"not null;type:varchar(512);"`
	Thumbnail    string `gorm:"type:varchar(512);"`
	BookWish     uint32 `gorm:"not null default 0"`
	BookDo       uint32 `gorm:"not null default 0"`
	BookCollect  uint32 `gorm:"not null default 0"`
	GameWish     uint32 `gorm:"not null default 0"`
	GameDo       uint32 `gorm:"not null default 0"`
	GameCollect  uint32 `gorm:"not null default 0"`
	MovieWish    uint32 `gorm:"not null default 0"`
	MovieDo      uint32 `gorm:"not null default 0"`
	MovieCollect uint32 `gorm:"not null default 0"`
	SongWish     uint32 `gorm:"not null default 0"`
	SongDo       uint32 `gorm:"not null default 0"`
	SongCollect  uint32 `gorm:"not null default 0"`
	SyncAt       time.Time
	CheckAt      time.Time
	RegisterAt   time.Time
	PublishAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (User) TableName() string {
	return "user"
}

func (user User) Show() *UserVO {
	return &UserVO{
		ID:           user.DoubanUid,
		Domain:       user.Domain,
		Name:         user.Name,
		Thumbnail:    user.Thumbnail,
		BookWish:     user.BookWish,
		BookDo:       user.BookDo,
		BookCollect:  user.BookCollect,
		GameWish:     user.GameWish,
		GameDo:       user.GameDo,
		GameCollect:  user.GameCollect,
		MovieWish:    user.MovieWish,
		MovieDo:      user.MovieDo,
		MovieCollect: user.MovieCollect,
		SongWish:     user.SongWish,
		SongDo:       user.SongDo,
		SongCollect:  user.SongCollect,
		SyncAt:       user.SyncAt.Unix(),
		CheckAt:      user.CheckAt.Unix(),
	}
}

type UserVO struct {
	ID           uint64 `json:"id"`
	Domain       string `json:"domain"`
	Name         string `json:"name"`
	Thumbnail    string `json:"thumbnail"`
	BookWish     uint32 `json:"book_wish"`
	BookDo       uint32 `json:"book_do"`
	BookCollect  uint32 `json:"book_collect"`
	GameWish     uint32 `json:"game_wish"`
	GameDo       uint32 `json:"game_do"`
	GameCollect  uint32 `json:"game_collect"`
	MovieWish    uint32 `json:"movie_wish"`
	MovieDo      uint32 `json:"movie_do"`
	MovieCollect uint32 `json:"movie_collect"`
	SongWish     uint32 `json:"song_wish"`
	SongDo       uint32 `json:"song_do"`
	SongCollect  uint32 `json:"song_collect"`
	SyncAt       int64  `json:"sync_at"`
	CheckAt      int64  `json:"check_at"`
}
