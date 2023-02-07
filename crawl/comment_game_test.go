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
	logrus.Info(util.ToJson(*comments))
	logrus.Info(util.ToJson(*games))
	logrus.Info(total)
	logrus.Info(next)
}
