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
func GetMovieDetail(doubanId uint64) *model.Movie {
	movie := &model.Movie{}
	common.Db.Where("douban_id = ? ", doubanId).Find(movie)
	return movie
}

func ListMovieBrief(doubanIds *[]uint64) *[]model.Movie {
	var movies *[]model.Movie
	common.Db.Omit("intro").Where("douban_id IN ? ", *doubanIds).Find(&movies)
	return movies
}
