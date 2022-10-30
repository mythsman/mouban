package crawl

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/consts"
	"mouban/model"
	"strings"
)

func CommentMovie(doubanUid int) (*model.User, []*model.Comment, error) {
	scrollMovie(doubanUid, "do")
	scrollMovie(doubanUid, "wish")
	scrollMovie(doubanUid, "collect")
	return nil, nil, nil
}

func scrollMovie(doubanUid int, action string) ([]*model.Comment, *int, *string, error) {
	body, err := Get(fmt.Sprintf(consts.MovieCommentUrl, doubanUid, action))
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