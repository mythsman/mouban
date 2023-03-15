package crawl

import (
	"github.com/antchfx/htmlquery"
	"mouban/consts"
	"mouban/util"
	"strings"
)

func LatestBook() *[]uint64 {
	body, _, err := Get(consts.BookLatestUrl, DiscoverLimiter)
	if err != nil {
		panic(err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))

	if err != nil {
		panic(err)
	}

	return util.ParseNewItems(doc, consts.TypeBook)
}

func LatestSong() *[]uint64 {
	body, _, err := Get(consts.SongLatestUrl, DiscoverLimiter)
	if err != nil {
		panic(err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))

	if err != nil {
		panic(err)
	}

	return util.ParseNewItems(doc, consts.TypeSong)
}

func LatestMovie() *[]uint64 {
	body, _, err := Get(consts.MovieLatestUrl, DiscoverLimiter)
	if err != nil {
		panic(err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))

	if err != nil {
		panic(err)
	}

	return util.ParseNewItems(doc, consts.TypeMovie)
}

func LatestGame() *[]uint64 {
	body, _, err := Get(consts.GameLatestUrl, DiscoverLimiter)
	if err != nil {
		panic(err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))

	if err != nil {
		panic(err)
	}

	return util.ParseNewItems(doc, consts.TypeGame)
}
