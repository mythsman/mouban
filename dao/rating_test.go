package dao

import (
	"mouban/consts"
	"mouban/model"
	"testing"
)

func TestUpsertRating(t *testing.T) {
	rating := &model.Rating{
		Total:  100,
		Rating: 3.3,
		Star5:  5.2,
		Star4:  3.4,
		Star3:  3.4,
		Star2:  2.4,
		Star1:  1.4,
		Status: consts.RatingNormal,
	}
	UpsertRating(rating)
}
