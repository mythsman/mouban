package model

import (
	"time"
)

type Storage struct {
	ID        uint64
	Source    string `gorm:"type:varchar(256);not null;uniqueIndex;"`
	Target    string `gorm:"type:varchar(256);not null;index;"`
	Md5       string `gorm:"type:varchar(64);not null;index;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Storage) TableName() string {
	return "storage"
}
