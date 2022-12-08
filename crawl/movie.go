package crawl

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strconv"
	"strings"
)

func Movie(doubanId uint64) (*model.Movie, *model.Rating, *[]string, *[]uint64, error) {
	body, _, err := Get(fmt.Sprintf(consts.MovieDetailUrl, doubanId), MovieLimiter)
	if err != nil {
		panic(err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		panic(err)
	}

	tt := htmlquery.FindOne(doc, "//head//title")
	if tt == nil {
		panic("title is nil for " + strconv.FormatUint(doubanId, 10) + ", html: {}" + htmlquery.OutputHTML(doc, true))
	}
	t := htmlquery.InnerText(tt)
	if strings.TrimSpace(t) == "页面不存在" || strings.TrimSpace(t) == "条目不存在" {
		return nil, nil, nil, nil, errors.New(strings.TrimSpace(t))
	}

	title := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//meta[@property='og:title']"), "content")
	thumbnail := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//a[@class='nbg']/img"), "src")

	intro := ""
	allHiddenIntro := htmlquery.FindOne(doc, "//div[@id='link-report-intra']/span[@class='all hidden']")
	if allHiddenIntro != nil {
		intro = util.TrimParagraph(htmlquery.InnerText(allHiddenIntro))
	} else {
		shortIntro := htmlquery.FindOne(doc, "//div[@id='link-report-intra']/span[@property='v:summary']")
		if shortIntro != nil {
			intro = util.TrimParagraph(htmlquery.InnerText(shortIntro))
		}
	}

	data := util.TrimInfo(htmlquery.OutputHTML(htmlquery.FindOne(doc, "//div[@id='info']"), false))

	director := strings.TrimSpace(data["编剧"])
	actor := strings.TrimSpace(data["主演"])
	writer := strings.TrimSpace(data["编剧"])
	site := strings.TrimSpace(data["官方网站"])
	style := strings.TrimSpace(data["类型"])
	country := strings.TrimSpace(data["制片国家/地区"])
	language := strings.TrimSpace(data["语言"])
	duration := uint64(0)
	if data["片长"] != "" {
		duration = util.ParseNumber(data["片长"]) * 60
	} else if data["单集片长"] != "" {
		duration = util.ParseNumber(data["单集片长"]) * 60
	}
	alias := strings.TrimSpace(data["又名"])
	imdb := strings.TrimSpace(data["IMDb"])
	episode := util.ParseNumber(data["集数"])
	releaseDate := strings.TrimSpace(data["上映日期"])

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
		PublishDate: releaseDate,
		Episode:     uint32(episode),
		Duration:    uint32(duration),
		Alias:       alias,
		IMDb:        imdb,
		Intro:       intro,
		Thumbnail:   thumbnail,
	}

	rating := Rating(htmlquery.FindOne(doc, "//div[@id='interest_sectl']"))
	rating.DoubanId = doubanId
	rating.Type = consts.TypeMovie.Code

	newUsers := util.ParseNewUsers(doc)
	newItems := util.ParseNewItems(doc, consts.TypeGame)

	return movie, rating, newUsers, newItems, nil
}
