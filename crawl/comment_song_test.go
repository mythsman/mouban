package crawl

import (
	"fmt"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func Test_scrollSong(t *testing.T) {
	comments, songs, total := scrollAllSong(43001468, consts.ActionCollect)
	fmt.Println(util.ToJson(*comments))
	fmt.Println(util.ToJson(*songs))
	fmt.Println(total)
}
