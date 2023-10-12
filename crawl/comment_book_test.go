package crawl

import (
	"mouban/consts"
	"mouban/util"
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_scrollBook(t *testing.T) {
	comments, books, total, next, err := scrollBook(214963638, "", consts.ActionCollect)
	if err != nil {
		return
	}
	logrus.Infoln(util.ToJson(*comments))
	logrus.Infoln(util.ToJson(*books))
	logrus.Infoln(total)
	logrus.Infoln(next)
}
