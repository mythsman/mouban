package model

import (
	"time"
)

type Storage struct {
	ID        uint64
	Source    uint64 `gorm:"not null;uniqueIndex;"`
	Target    uint64 `gorm:"not null;index;"`
	Md5       string `gorm:"not null;index;"`
	Extra     string `gorm:"type:varchar(512)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Storage) TableName() string {
	return "storage"
}
