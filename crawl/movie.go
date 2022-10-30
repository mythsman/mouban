package crawl

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/consts"
	"mouban/model"
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
	list := htmlquery.Find(doc, "//a[@href]")
	for i := range list {
		fmt.Println(list[i])
	}

	return nil, nil, nil
}
