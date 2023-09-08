package crawl

import (
	"github.com/sirupsen/logrus"
	"mouban/consts"
	"mouban/util"
	"testing"
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
