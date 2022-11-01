package dao

import (
	"mouban/common"
	"mouban/model"
)

func UpsertMovie(movie *model.Movie) {
	if common.Db.Where("douban_id = ? ", movie.DoubanId).Updates(movie).RowsAffected == 0 {
		common.Db.Create(movie)
	}
}
