package model

import (
	"time"
)

type Storage struct {
	ID        uint64
	Source    string `gorm:"not null;uniqueIndex;"`
	Target    string `gorm:"not null;index;"`
	Md5       string `gorm:"not null;index;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Storage) TableName() string {
	return "storage"
}
