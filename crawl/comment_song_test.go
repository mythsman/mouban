package crawl

import (
	"github.com/sirupsen/logrus"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func Test_scrollSong(t *testing.T) {
	comments, songs, total := scrollAllSong(43001468, consts.ActionCollect)
	logrus.Infoln(util.ToJson(*comments))
	logrus.Infoln(util.ToJson(*songs))
	logrus.Infoln(total)
}
