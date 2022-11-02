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

	intro := ""
	allHiddenIntro := htmlquery.FindOne(doc, "//div[@id='link-report-intra']/span[@class='all hidden']")
	if allHiddenIntro != nil {
		intro = util.TrimParagraph(htmlquery.InnerText(allHiddenIntro))
	} else {
		shortIntro := htmlquery.FindOne(doc, "//div[@id='link-report-intra']/span[@property='v:summary']")
		intro = util.TrimParagraph(htmlquery.InnerText(shortIntro))
	}

	data := util.TrimInfo(htmlquery.OutputHTML(htmlquery.FindOne(doc, "//div[@id='info']"), false))

	director := data["编剧"]
	actor := data["主演"]
	writer := data["编剧"]
	site := data["官方网站"]
	style := data["类型"]
	country := data["制片国家/地区"]
	language := data["语言"]
	duration := uint64(0)
	if data["片长"] != "" {
		duration = util.ParseNumber(data["片长"]) * 60
	} else if data["单集片长"] != "" {
		duration = util.ParseNumber(data["单集片长"]) * 60
	}
	alias := data["又名"]
	imdb := data["Imdb"]
	episode := util.ParseNumber(data["集数"])
	releaseData := data["上映日期"]

	movie := &model.Movie{
		DoubanId:    doubanId,
		Title:       title,
		Director:    director,
		Writer:      writer,
		Actor:       actor,
		Style:       style,
		Site:        site,
		Country:     country,
		Language:    language,
		PublishDate: releaseData,
		Episode:     uint32(episode),
		Duration:    uint32(duration),
		Alias:       alias,
		IMDb:        imdb,
		Intro:       intro,
		Thumbnail:   thumbnail,
	}

	rating := Rating(htmlquery.FindOne(doc, "//div[@id='interest_sectl']"))
	rating.DoubanId = doubanId
	rating.Type = consts.TypeMovie

	return movie, rating, nil
}
