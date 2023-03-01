package dao

import (
	"mouban/common"
	"mouban/model"
)

func UpsertRating(rating *model.Rating) {
	data := &model.Rating{}
	common.Db.Where("douban_id = ? AND type = ?", rating.DoubanId, rating.Type).Assign(rating).FirstOrCreate(data)
}

func GetRating(doubanId uint64, t uint8) *model.Rating {
	rating := &model.Rating{}
	common.Db.Where("douban_id = ? AND type = ?", doubanId, t).Find(rating)
	if rating.ID == 0 {
		return nil
	}
	return rating
}

func ListRating(doubanIds *[]uint64, t uint8) *[]model.Rating {
	var rating *[]model.Rating
	common.Db.Where("douban_id IN ? AND type = ?", *doubanIds, t).Find(&rating)
	return rating
}
