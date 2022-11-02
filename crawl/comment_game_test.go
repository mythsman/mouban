package crawl

import (
	"fmt"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func Test_scrollGame(t *testing.T) {
	comments, games, total, next, err := scrollGame(162448367, "", consts.ActionCollect)
	if err != nil {
		return
	}
	fmt.Println(util.ToJson(*comments))
	fmt.Println(util.ToJson(*games))
	fmt.Println(total)
	fmt.Println(next)
}
