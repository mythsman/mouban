package model

import (
	"time"
)

type Comment struct {
	ID        uint64
	DoubanUid uint64 `gorm:"not null;uniqueIndex:uk_comment;index:idx_search;priority:1"`
	DoubanId  uint64 `gorm:"not null;uniqueIndex:uk_comment;priority:2"`
	Type      uint8  `gorm:"not null;uniqueIndex:uk_comment;index:idx_search;priority:3"`
	Rate      uint8
	Label     string    `gorm:"type:varchar(512)"`
	Comment   string    `gorm:"type:mediumtext"`
	Action    uint8     `gorm:"not null;index:idx_search;priority:4"`
	MarkDate  time.Time `gorm:"not null;index:idx_search;priority:5"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Comment) TableName() string {
	return "comment"
}

func (comment Comment) Show(item interface{}) *CommentVO {
	return &CommentVO{
		Item:     item,
		Rate:     comment.Rate,
		Label:    comment.Label,
		Comment:  comment.Comment,
		Action:   comment.Action,
		MarkDate: comment.MarkDate,
	}
}

type CommentVO struct {
	Item     interface{} `json:"item"`
	Rate     uint8       `json:"rate"`
	Label    string      `json:"label"`
	Comment  string      `json:"comment"`
	Action   uint8       `json:"action"`
	MarkDate time.Time   `json:"mark_date"`
}
