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

func Song(doubanId uint64) (*model.Song, *model.Rating, *[]string, *[]uint64, error) {
	body, _, err := Get(fmt.Sprintf(consts.SongDetailUrl, doubanId), ItemLimiter)
	if err != nil {
		if strings.Contains(err.Error(), "too many redirects") {
			return nil, nil, nil, nil, err
		}
		panic(err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		panic(err)
	}

	tt := htmlquery.FindOne(doc, "//head//title")
	if tt == nil {
		panic("title is nil for " + strconv.FormatUint(doubanId, 10) + ", html: " + htmlquery.OutputHTML(doc, true))
	}
	t := htmlquery.InnerText(tt)
	if strings.TrimSpace(t) == "页面不存在" || strings.TrimSpace(t) == "条目不存在" {
		return nil, nil, nil, nil, errors.New(strings.TrimSpace(t))
	}

	title := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//meta[@property='og:title']"), "content")
	thumbnailNode := htmlquery.FindOne(doc, "//a[@class='nbg']/img")
	if thumbnailNode == nil {
		thumbnailNode = htmlquery.FindOne(doc, "//a[@class='nbgnbg']/img")
	}

	thumbnail := htmlquery.SelectAttr(thumbnailNode, "src")

	intro := ""
	allHiddenIntroNode := htmlquery.FindOne(doc, "//div[@id='link-report']/span[@class='all hidden']")
	if allHiddenIntroNode != nil {
		intro = util.TrimParagraph(htmlquery.OutputHTML(allHiddenIntroNode, false))
	} else {
		shortIntroNode := htmlquery.FindOne(doc, "//div[@id='link-report']/span[@property='v:summary']")
		if shortIntroNode != nil {
			intro = util.TrimParagraph(htmlquery.OutputHTML(shortIntroNode, false))
		}
	}

	trackList := ""
	trackListNode := htmlquery.FindOne(doc, "//div[@class='track-list']/div[@class='indent']/div")
	if trackListNode != nil {
		trackList = util.TrimParagraph(htmlquery.OutputHTML(trackListNode, false))
	}

	data := util.TrimInfo(htmlquery.OutputHTML(htmlquery.FindOne(doc, "//div[@id='info']"), false))

	alias := strings.TrimSpace(data["又名"])
	musician := strings.TrimSpace(data["表演者"])
	albumType := strings.TrimSpace(data["专辑类型"])
	genre := strings.TrimSpace(data["流派"])
	media := strings.TrimSpace(data["介质"])
	barcode := strings.TrimSpace(data["条形码"])
	publisher := strings.TrimSpace(data["出版者"])
	publishDate := strings.TrimSpace(data["发行时间"])
	ISRC := strings.TrimSpace(data["ISRC(中国)"])
	albumCount := util.ParseNumber(data["唱片数"])

	movie := &model.Song{
		DoubanId:    doubanId,
		Title:       title,
		Alias:       alias,
		Musician:    musician,
		AlbumType:   albumType,
		Genre:       genre,
		Media:       media,
		Barcode:     barcode,
		Publisher:   publisher,
		PublishDate: publishDate,
		ISRC:        ISRC,
		AlbumCount:  uint32(albumCount),
		Intro:       intro,
		TrackList:   trackList,
		Thumbnail:   thumbnail,
	}

	rating := Rating(htmlquery.FindOne(doc, "//div[@id='interest_sectl']"))
	rating.DoubanId = doubanId
	rating.Type = consts.TypeSong.Code

	newUsers := util.ParseNewUsers(doc)
	newItems := util.ParseNewItems(doc, consts.TypeSong)

	return movie, rating, newUsers, newItems, nil
}
