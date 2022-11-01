package dao

import (
	"mouban/model"
	"testing"
)

func TestUpsertGame(t *testing.T) {
	game := &model.Game{
		DoubanId:    2342,
		Title:       "title",
		Platform:    "平台",
		Genre:       "类型",
		Alias:       "别名",
		Developer:   "开发商",
		Publisher:   "发行商",
		PublishDate: "发行日期",
		Intro:       "intro",
		Thumbnail:   "thumbnail",
	}
	UpsertGame(game)
}
