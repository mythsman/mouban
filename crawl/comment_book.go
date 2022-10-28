package crawl

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/common"
	"mouban/model"
	"strings"
)

func CommentBook(doubanUid int) (*model.User, []*model.Comment, error) {
	scrollBook(doubanUid, "do")
	scrollBook(doubanUid, "wish")
	scrollBook(doubanUid, "collect")
	return nil, nil, nil
}

func scrollBook(doubanUid int, action string) ([]*model.Comment, *int, *string, error) {
	body, err := Get(fmt.Sprintf(common.BookCommentUrl, doubanUid, action))
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
