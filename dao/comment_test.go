package dao

import (
	"mouban/consts"
	"mouban/model"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestUpsertComment(t *testing.T) {
	comment := &model.Comment{
		DoubanUid: 11,
		DoubanId:  23,
		Type:      consts.TypeBook.Code,
		Rate:      3,
		Label:     "tags",
		Comment:   "shortComment",
		Action:    &consts.ActionWish.Code,
		MarkDate:  time.Now(),
	}
	UpsertComment(comment)
}

func TestGetCommentIds(t *testing.T) {
	data := GetCommentIds(11, consts.TypeMovie.Code)
	logrus.Infoln("data", *data)
}
