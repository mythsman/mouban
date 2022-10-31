package crawl

import (
	"fmt"
	"mouban/consts"
	"testing"
)

func Test_scrollMovie(t *testing.T) {
	movie, total, next, err := scrollMovie(162448367, "", consts.ActionCollect)
	if err != nil {
		return
	}
	for i := range movie {
		fmt.Println(movie[i])
	}
	fmt.Println(total)
	fmt.Println(next)
}
