package crawl

import (
	"github.com/sirupsen/logrus"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func Test_scrollSong(t *testing.T) {
	comments, songs, total := scrollAllSong(43001468, consts.ActionCollect)
	logrus.Info(util.ToJson(*comments))
	logrus.Info(util.ToJson(*songs))
	logrus.Info(total)
}
