package crawl

import (
	"github.com/sirupsen/logrus"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func Test_scrollGame(t *testing.T) {
	comments, games, total, next, err := scrollGame(162448367, "", consts.ActionCollect)
	if err != nil {
		return
	}
	logrus.Infoln(util.ToJson(*comments))
	logrus.Infoln(util.ToJson(*games))
	logrus.Infoln(total)
	logrus.Infoln(next)
}
