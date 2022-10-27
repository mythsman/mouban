package do

import (
	"time"
)

type Comment struct {
	DoubanId uint64    `json:"douban_id"`
	Type     uint8     `json:"type"`
	Rate     uint8     `json:"rate"`
	Label    string    `json:"label"`
	Comment  string    `json:"comment"`
	Status   uint8     `json:"status"`
	MarkDate time.Time `json:"mark_date"`
}
