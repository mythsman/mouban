package crawl

import (
	"github.com/sirupsen/logrus"
	"mouban/consts"
	"mouban/util"
	"testing"
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
	user, comments, books, err := CommentBook(162448367)
	if err != nil {
		return
	}
	logrus.Infoln(util.ToJson(*user))
	logrus.Infoln(util.ToJson(*comments))
	logrus.Infoln(util.ToJson(*books))

}
