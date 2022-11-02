package crawl

import (
	"fmt"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func Test_scrollBook(t *testing.T) {
	comments, books, total, next, err := scrollBook(162448367, "", consts.ActionCollect)
	if err != nil {
		return
	}
	fmt.Println(util.ToJson(*comments))
	fmt.Println(util.ToJson(*books))
	fmt.Println(total)
	fmt.Println(next)
}

func Test_CommentBook(t *testing.T) {
	user, comments, books, err := CommentBook(162448367)
	if err != nil {
		return
	}
	fmt.Println(util.ToJson(*user))
	fmt.Println(util.ToJson(*comments))
	fmt.Println(util.ToJson(*books))

}
