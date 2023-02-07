package dao

import (
	"github.com/sirupsen/logrus"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"testing"
)

func TestUpsertRating(t *testing.T) {
	rating := &model.Rating{
		DoubanId: 123,
		Type:     consts.TypeMovie.Code,
		Total:    100,
		Rating:   3.3,
		Star5:    5.2,
		Star4:    3.4,
		Star3:    3.4,
		Star2:    2.4,
		Star1:    1.4,
		Status:   consts.RatingNormal,
	}
	UpsertRating(rating)
}

func TestGetRating(t *testing.T) {
	rating := GetRating(123, consts.TypeMovie.Code)
	logrus.Infoln(util.ToJson(rating))
}

func TestListRating(t *testing.T) {
	ratings := ListRating(&[]uint64{123, 123}, consts.TypeMovie.Code)
	logrus.Infoln(util.ToJson(ratings))
}
