package crawl

import (
	"mouban/consts"
	"mouban/log"
	"mouban/util"
	"testing"
)

func Test_scrollGame(t *testing.T) {
	comments, games, total, next, err := scrollGame(162448367, "", consts.ActionCollect)
	if err != nil {
		return
	}
	log.Info(util.ToJson(*comments))
	log.Info(util.ToJson(*games))
	log.Info(total)
	log.Info(next)
}
