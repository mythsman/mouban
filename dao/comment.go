package dao

import (
	"github.com/sirupsen/logrus"
	"mouban/common"
	"mouban/consts"
	"mouban/model"
)

func HideComment(doubanUid uint64, t uint8, doubanId uint64) {
	logrus.Infoln("hide comment for", doubanUid, "type", t, "at", doubanId)
	common.Db.Model(&model.Comment{}).
		Where("douban_uid = ? AND type = ? AND douban_id = ?", doubanUid, t, doubanId).
		Update("action", consts.ActionHide.Code)
}

func GetCommentIds(doubanUid uint64, t uint8) *[]uint64 {
	var doubanIds []uint64
	common.Db.Model(&model.Comment{}).Where("douban_uid = ? AND type = ?", doubanUid, t).Select("douban_id").Find(&doubanIds)
	return &doubanIds
}

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
