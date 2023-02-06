package crawl

import (
	"github.com/sirupsen/logrus"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func Test_scrollSong(t *testing.T) {
	comments, songs, total := scrollAllSong(43001468, consts.ActionCollect)
	log.Info(util.ToJson(*comments))
	log.Info(util.ToJson(*songs))
	log.Info(total)
}
