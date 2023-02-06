package crawl

import (
	"mouban/consts"
	"mouban/log"
	"mouban/util"
	"testing"
)

func Test_scrollBook(t *testing.T) {
	comments, books, total, next, err := scrollBook(162448367, "", consts.ActionCollect)
	if err != nil {
		return
	}
	log.Info(util.ToJson(*comments))
	log.Info(util.ToJson(*books))
	log.Info(total)
	log.Info(next)
}

func Test_CommentBook(t *testing.T) {
	user, comments, books, err := CommentBook(162448367)
	if err != nil {
		return
	}
	log.Info(util.ToJson(*user))
	log.Info(util.ToJson(*comments))
	log.Info(util.ToJson(*books))

}
