package crawl

import (
	"fmt"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func Test_scrollBook(t *testing.T) {
	book, total, next, err := scrollBook(162448367, "", consts.ActionCollect)
	if err != nil {
		return
	}
	for i := range book {
		fmt.Println(util.ToJson(book[i]))
	}
	fmt.Println(total)
	fmt.Println(next)
}

func Test_CommentBook(t *testing.T) {
	user, comments, err := CommentBook(162448367)
	if err != nil {
		return
	}
	fmt.Println(util.ToJson(user))
	for _, c := range comments {
		fmt.Println(util.ToJson(c))
	}
}
