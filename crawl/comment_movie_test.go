package crawl

import (
	"fmt"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func Test_scrollMovie(t *testing.T) {
	comments, movies, total, next, err := scrollMovie(162448367, "", consts.ActionCollect)
	if err != nil {
		return
	}
	fmt.Println(util.ToJson(*comments))
	fmt.Println(util.ToJson(*movies))

	fmt.Println(total)
	fmt.Println(next)
}
