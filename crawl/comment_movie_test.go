package crawl

import (
	"fmt"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func Test_scrollMovie(t *testing.T) {
	movie, total, next, err := scrollMovie(162448367, "", consts.ActionCollect)
	if err != nil {
		return
	}
	for i := range movie {
		fmt.Println(util.ToJson(movie[i]))
	}
	fmt.Println(total)
	fmt.Println(next)
}
