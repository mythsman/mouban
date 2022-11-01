package dao

import (
	"mouban/model"
	"testing"
)

func TestUpsertMovie(t *testing.T) {

	movie := &model.Movie{
		DoubanId:    22,
		Title:       "title",
		Director:    "director",
		Writer:      "writer",
		Actor:       "actor",
		Style:       "style",
		Site:        "site",
		Country:     "country",
		Language:    "language",
		PublishDate: "publishDate",
		Episode:     1,
		Duration:    2,
		Alias:       "alias",
		IMDb:        "imdb",
		Intro:       "intro",
		Thumbnail:   "thumbnail",
	}

	UpsertMovie(movie)
}
