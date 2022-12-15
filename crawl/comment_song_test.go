package crawl

import (
	"fmt"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func Test_scrollSong(t *testing.T) {
	comments, songs, total, next, err := scrollSong(43001468, "https://music.douban.com/people/43001468/do?start=0&sort=time&rating=all&filter=all&mode=grid", consts.ActionCollect)
	if err != nil {
		return
	}
	fmt.Println(util.ToJson(*comments))
	fmt.Println(util.ToJson(*songs))

	fmt.Println(total)
	fmt.Println(next)
}
