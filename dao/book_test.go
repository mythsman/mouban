package dao

import (
	"mouban/model"
	"testing"
)

func TestUpsertBook(t *testing.T) {
	book := &model.Book{
		DoubanId:   11,
		Title:      "title",
		Subtitle:   "subtitle",
		Orititle:   "orititle",
		Author:     "author",
		Translator: "translator",
		Press:      "press",
		Producer:   "producer",
		Serial:     "serial",
		PublishAt:  "publishAt",
		ISBN:       "isbn",
		Framing:    "framing",
		Page:       0,
		Price:      0,
		Intro:      "intro",
		Thumbnail:  "thumbnail",
	}
	UpsertBook(book)
}
