package crawl

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/common"
	"mouban/model"
	"strings"
)

func CommentGame(doubanUid int) (*model.User, []*model.Comment, error) {
	scrollGame(doubanUid, "do")
	scrollGame(doubanUid, "wish")
	scrollGame(doubanUid, "collect")
	return nil, nil, nil
}

func scrollGame(doubanUid int, action string) ([]*model.Comment, *int, *string, error) {
	body, err := Get(fmt.Sprintf(common.GameCommentUrl, doubanUid, action))
	if err != nil {
		return nil, nil, nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, nil, nil, err
	}
	list := htmlquery.Find(doc, "//a[@href]")
	for i := range list {
		fmt.Println(list[i])
	}

	return nil, nil, nil, err

}
