package crawl

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/model"
	"mouban/util"
	"strings"
)

func Game(doubanId int) (*model.Game, *model.Rating, error) {
	body, err := Get(fmt.Sprintf(util.GameDetailUrl, doubanId))
	if err != nil {
		return nil, nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, nil, err
	}
	list := htmlquery.Find(doc, "//a[@href]")
	for i := range list {
		fmt.Println(list[i])
	}

	return nil, nil, nil
}
