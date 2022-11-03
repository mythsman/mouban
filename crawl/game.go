package crawl

import (
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strings"
)

func Game(doubanId uint64) (*model.Game, *model.Rating, error) {
	body, err := Get(fmt.Sprintf(consts.GameDetailUrl, doubanId))
	if err != nil {
		return nil, nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, nil, err
	}

	t := htmlquery.InnerText(htmlquery.FindOne(doc, "//head//title"))
	if strings.TrimSpace(t) == "页面不存在" {
		return nil, nil, errors.New("页面不存在")
	}

	title := htmlquery.InnerText(htmlquery.FindOne(doc, "//div[@id='content']/h1"))
	thumbnail := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//div[@class='pic']/img"), "src")
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
		Platform:    data["平台"],
		Genre:       data["类型"],
		Alias:       data["别名"],
		Developer:   data["开发商"],
		Publisher:   data["发行商"],
		PublishDate: data["发行日期"],
		Intro:       intro,
		Thumbnail:   thumbnail,
	}

	rating := Rating(htmlquery.FindOne(doc, "//div[@id='interest_sectl']"))
	rating.DoubanId = doubanId
	rating.Type = consts.TypeGame
	return game, rating, nil
}
