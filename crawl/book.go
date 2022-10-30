package crawl

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strings"
)

func Book(doubanId uint64) (*model.Book, *model.Rating, error) {
	body, err := Get(fmt.Sprintf(consts.BookDetailUrl, doubanId))
	if err != nil {
		return nil, nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, nil, err
	}
	title := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//meta[@property='og:title']"), "content")
	thumbnail := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//a[@class='nbg']/img"), "src")
	intros := htmlquery.Find(doc, "//div[@class='intro']")
	intro := strings.TrimSpace(htmlquery.InnerText(intros[len(intros)-1]))

	data := util.TrimInfo(htmlquery.OutputHTML(htmlquery.FindOne(doc, "//div[@id='info']"), false))

	isbn := data["ISBN"]
	subtitle := data["副标题"]
	orititle := data["原作名"]
	author := data["作者"]
	press := data["出版社"]
	producer := data["出品方"]
	translator := data["译者"]
	serial := data["丛书"]
	publishAt := data["出版年"]
	page := uint32(util.ParseNumber(data["页数"]))
	price := uint32(util.ParseFloat(data["定价"]) * 100)

	book := &model.Book{
		DoubanId:   doubanId,
		Title:      title,
		Subtitle:   subtitle,
		Orititle:   orititle,
		Author:     author,
		Translator: translator,
		Press:      press,
		Producer:   producer,
		Serial:     serial,
		PublishAt:  publishAt,
		ISBN:       isbn,
		Page:       page,
		Price:      price,
		Intro:      intro,
		Thumbnail:  thumbnail,
	}

	rating := Rating(htmlquery.FindOne(doc, "//div[@id='interest_sectl']"))

	return book, rating, nil
}
