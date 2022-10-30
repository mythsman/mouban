package crawl

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strings"
)

func Movie(doubanId uint64) (*model.Movie, *model.Rating, error) {
	body, err := Get(fmt.Sprintf(consts.MovieDetailUrl, doubanId))
	if err != nil {
		return nil, nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, nil, err
	}
	title := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//meta[@property='og:title']"), "content")
	thumbnail := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//a[@class='nbg']/img"), "src")
	intro := strings.TrimSpace(htmlquery.InnerText(htmlquery.FindOne(doc, "//div[@id='link-report-intra']")))

	data := util.TrimInfo(htmlquery.OutputHTML(htmlquery.FindOne(doc, "//div[@id='info']"), false))
	fmt.Println(data)

	movie := &model.Movie{
		DoubanId:    doubanId,
		Title:       title,
		Director:    "",
		Writer:      "",
		Actor:       "",
		Style:       "",
		Site:        "",
		Country:     "",
		Language:    "",
		ReleaseDate: "",
		Season:      0,
		Episode:     0,
		Duration:    0,
		Alias:       "",
		IMDb:        "",
		Intro:       intro,
		Thumbnail:   thumbnail,
	}

	rating := Rating(htmlquery.FindOne(doc, "//div[@id='interest_sectl']"))

	return movie, rating, nil
}
