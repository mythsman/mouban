package dao

import (
	"mouban/internal/common"
	"mouban/internal/model"

	"github.com/sirupsen/logrus"
)

func UpsertRating(rating *model.Rating) {
	logrus.Infoln("upsert rating", rating.DoubanId, rating.Type)
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
