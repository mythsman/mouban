package crawl

import (
	"mouban/internal/consts"
	"mouban/internal/util"
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_scrollGame(t *testing.T) {
	comments, games, total, next, err := scrollGame(162448367, "", consts.ActionCollect)
	if err != nil {
		t.Fatalf("scrollGame failed: %v", err)
	}
	logrus.Infoln(util.ToJson(*comments))
	logrus.Infoln(util.ToJson(*games))
	logrus.Infoln(total)
	logrus.Infoln(next)
}
