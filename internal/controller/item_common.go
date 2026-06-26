package controller

import (
	"strings"
	"time"

	"mouban/internal/consts"
)

func formatTimeCN(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(time.Local).Format("2006-01-02 15:04:05")
}

func parseItemType(typeName string) (uint8, string, string, string) {
	switch strings.ToLower(strings.TrimSpace(typeName)) {
	case "book":
		return consts.TypeBook.Code, "book", "图书", "https://book.douban.com/subject/"
	case "movie":
		return consts.TypeMovie.Code, "movie", "电影", "https://movie.douban.com/subject/"
	case "game":
		return consts.TypeGame.Code, "game", "游戏", "https://www.douban.com/game/"
	case "song":
		return consts.TypeSong.Code, "song", "音乐", "https://music.douban.com/subject/"
	default:
		return 0, "", "", ""
	}
}
