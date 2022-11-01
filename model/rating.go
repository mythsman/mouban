package model

import (
	"time"
)

type Rating struct {
	ID        uint64
	Type      uint8  `gorm:"not null;uniqueIndex:uk_unique_id"`
	DoubanId  uint64 `gorm:"not null;uniqueIndex:uk_unique_id"`
	Total     uint32
	Rating    float32
	Star5     float32
	Star4     float32
	Star3     float32
	Star2     float32
	Star1     float32
	Status    uint8
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Rating) TableName() string {
	return "rating"
}

type RatingVO struct {
	Total  uint32  `json:"total"`
	Rating float32 `json:"rating"`
	Star5  float32 `json:"star5"`
	Star4  float32 `json:"star4"`
	Star3  float32 `json:"star3"`
	Star2  float32 `json:"star2"`
	Star1  float32 `json:"star1"`
	Status string  `json:"status"`
}
