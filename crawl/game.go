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

func Game(doubanId uint64) (*model.Game, *model.Rating, *[]string, *[]uint64, error) {
	body, _, err := Get(fmt.Sprintf(consts.GameDetailUrl, doubanId), ItemLimiter)
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

	title := htmlquery.InnerText(htmlquery.FindOne(doc, "//div[@id='content']/h1"))
	thumbnailNode := htmlquery.FindOne(doc, "//div[@class='pic']//img")
	thumbnail := htmlquery.SelectAttr(thumbnailNode, "src")
	intro := util.TrimParagraph(htmlquery.InnerText(htmlquery.FindOne(doc, "//div[@id='link-report']/p")))

	data := make(map[string]string)
	labels := htmlquery.Find(doc, "//dl[@class='game-attr']/dt")
	values := htmlquery.Find(doc, "//dl[@class='game-attr']/dd")

	for i := range labels {
		data[strings.Trim(util.TrimLine(htmlquery.InnerText(labels[i])), ":")] = util.TrimLine(htmlquery.InnerText(values[i]))
	}

	game := &model.Game{
		DoubanId:    doubanId,
		Title:       title,
		Platform:    strings.TrimSpace(data["平台"]),
		Genre:       strings.TrimSpace(data["类型"]),
		Alias:       strings.TrimSpace(data["别名"]),
		Developer:   strings.TrimSpace(data["开发商"]),
		Publisher:   strings.TrimSpace(data["发行商"]),
		PublishDate: strings.TrimSpace(data["发行日期"]),
		Intro:       intro,
		Thumbnail:   thumbnail,
	}

	rating := Rating(htmlquery.FindOne(doc, "//div[@id='interest_sectl']"))
	rating.DoubanId = doubanId
	rating.Type = consts.TypeGame.Code

	newUsers := util.ParseNewUsers(doc)
	newItems := util.ParseNewItems(doc, consts.TypeGame)

	return game, rating, newUsers, newItems, nil
}
