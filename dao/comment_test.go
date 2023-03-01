package dao

import (
	"mouban/consts"
	"mouban/model"
	"testing"
	"time"
)

func TestUpsertComment(t *testing.T) {
	comment := &model.Comment{
		DoubanUid: 11,
		DoubanId:  22,
		Type:      consts.TypeMovie.Code,
		Rate:      3,
		Label:     "tags",
		Comment:   "shortComment",
		Action:    &consts.ActionWish.Code,
		MarkDate:  time.Now(),
	}
	UpsertComment(comment)
}
