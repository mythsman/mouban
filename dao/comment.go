package dao

import (
	"github.com/sirupsen/logrus"
	"mouban/common"
	"mouban/model"
)

func UpsertComment(comment *model.Comment) {
	logrus.Infoln("upsert comment", comment.DoubanId, comment.Type, "for", comment.DoubanUid)
	data := &model.Comment{}
	common.Db.Where("douban_id = ? AND douban_uid = ? AND type = ?", comment.DoubanId, comment.DoubanUid, comment.Type).Assign(comment).FirstOrCreate(data)
}

func GetComment(doubanId uint64, doubanUid uint64, t uint8) *model.Comment {
	comment := &model.Comment{}
	common.Db.Where("douban_id = ? AND douban_uid = ? AND type = ?", doubanId, doubanUid, t).Find(comment)
	if comment.ID == 0 {
		return nil
	}
	return comment
}

// SearchComment idx_search
func SearchComment(doubanUid uint64, t uint8, action uint8, offset int) *[]model.Comment {
	var comment *[]model.Comment
	common.Db.Where("douban_uid = ? AND type = ? AND action = ? ", doubanUid, t, action).
		Order("mark_date desc").
		Offset(offset).
		Find(&comment)
	return comment
}
