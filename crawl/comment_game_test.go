package crawl

import (
	"fmt"
	"mouban/consts"
	"testing"
)

func Test_scrollGame(t *testing.T) {
	game, total, next, err := scrollGame(162448367, "", consts.ActionCollect)
	if err != nil {
		return
	}
	for i := range game {
		fmt.Println(game[i])
	}
	fmt.Println(total)
	fmt.Println(next)
}
