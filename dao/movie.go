package dao

import (
	"mouban/common"
	"mouban/model"
)

func UpsertMovie(movie *model.Movie) {
	data := &model.Movie{}
	common.Db.Where("douban_id = ? ", movie.DoubanId).Assign(movie).FirstOrCreate(data)
}

func CreateMovieNx(movie *model.Movie) bool {
	data := &model.Movie{}
	result := common.Db.Where("douban_id = ? ", movie.DoubanId).Attrs(movie).FirstOrCreate(data)
	return result.RowsAffected > 0
}

func GetMovieDetail(doubanId uint64) *model.Movie {
	movie := &model.Movie{}
	common.Db.Where("douban_id = ? ", doubanId).Find(movie)
	if movie.ID == 0 {
		return nil
	}
	return movie
}

func ListMovieBrief(doubanIds *[]uint64) *[]model.Movie {
	var movies *[]model.Movie
	common.Db.Omit("intro").Where("douban_id IN ? ", *doubanIds).Find(&movies)
	return movies
}
