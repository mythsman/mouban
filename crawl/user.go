package crawl

import (
	"crypto/md5"
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/common"
	"mouban/model"
	"strings"
)

func UserFull(id string) (*model.User, error) {
	UserHash(id)
	bookOverview(id)
	movieOverview(id)
	musicOverview(id)
	gameOverview(id)
	return nil, nil
}

func UserHash(id string) (*string, error) {
	body, err := Get(fmt.Sprintf(common.UserRssUrl, id))
	if err != nil {
		return nil, err
	}
	data := []byte(*body)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)

	return &md5str, nil
}

func bookOverview(id string) (*model.User, error) {
	body, err := Get(fmt.Sprintf(common.BookOverviewUrl, id))
	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, err

	}
	list := htmlquery.Find(doc, "//a[@href]")
	for i := range list {
		fmt.Println(list[i])
	}
	return nil, err

}

func movieOverview(id string) (*model.User, error) {
	body, err := Get(fmt.Sprintf(common.MovieOverviewUrl, id))
	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, err

	}
	list := htmlquery.Find(doc, "//a[@href]")
	for i := range list {
		fmt.Println(list[i])
	}
	return nil, err

}

func gameOverview(id string) (*model.User, error) {
	body, err := Get(fmt.Sprintf(common.GameOverviewUrl, id))
	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, err
	}
	list := htmlquery.Find(doc, "//a[@href]")
	for i := range list {
		fmt.Println(list[i])
	}
	return nil, err

}

func musicOverview(id string) (*model.User, error) {
	body, err := Get(fmt.Sprintf(common.MusicOverviewUrl, id))
	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, err

	}
	list := htmlquery.Find(doc, "//a[@href]")
	for i := range list {
		fmt.Println(list[i])
	}
	return nil, err

}
