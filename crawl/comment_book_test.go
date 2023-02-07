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
	logrus.Info(util.ToJson(*comments))
	logrus.Info(util.ToJson(*books))
	logrus.Info(total)
	logrus.Info(next)
}

func Test_CommentBook(t *testing.T) {
	user, comments, books, err := CommentBook(162448367)
	if err != nil {
		return
	}
	logrus.Info(util.ToJson(*user))
	logrus.Info(util.ToJson(*comments))
	logrus.Info(util.ToJson(*books))

}
