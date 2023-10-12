package crawl

import (
	"mouban/consts"
	"mouban/util"
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_scrollMovie(t *testing.T) {
	comments, movies, total, next, err := scrollMovie(221941333, "", consts.ActionCollect)
	if err != nil {
		return
	}
	logrus.Infoln(util.ToJson(*comments))
	logrus.Infoln(util.ToJson(*movies))

	logrus.Infoln(total)
	logrus.Infoln(next)
}
