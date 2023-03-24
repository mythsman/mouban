package crawl

import (
	"github.com/sirupsen/logrus"
	"mouban/consts"
	"mouban/util"
	"testing"
	"time"
)

func Test_scrollBook(t *testing.T) {
	comments, books, total, next, err := scrollBook(162448367, "", consts.ActionCollect)
	if err != nil {
		return
	}
	logrus.Infoln(util.ToJson(*comments))
	logrus.Infoln(util.ToJson(*books))
	logrus.Infoln(total)
	logrus.Infoln(next)
}

func Test_CommentBook(t *testing.T) {
	comments, books, err := CommentBook(162448367, time.Unix(0, 0))
	if err != nil {
		return
	}
	logrus.Infoln(util.ToJson(*comments))
	logrus.Infoln(util.ToJson(*books))

}
