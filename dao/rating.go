package dao

import (
	"mouban/common"
	"mouban/model"
)

func UpsertRating(rating *model.Rating) {
	if common.Db.Where("douban_id = ? AND type = ?", rating.DoubanId, rating.Type).Updates(rating).RowsAffected == 0 {
		common.Db.Create(rating)
	}
}
