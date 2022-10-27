package do

import (
	"time"
)

type Game struct {
	ID          uint64
	DoubanId    uint64 `gorm:"not null;uniqueIndex"`
	Title       string `gorm:"not null;type:varchar(512)"`
	Platform    string `gorm:"type:varchar(512)"`
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
