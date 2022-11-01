package dao

import (
	"mouban/common"
	"mouban/model"
)

func UpsertComment(comment *model.Comment) {
	if common.Db.Where("douban_id = ? AND douban_uid = ? AND type = ?", comment.DoubanId, comment.DoubanUid, comment.Type).
		Updates(comment).RowsAffected == 0 {
		common.Db.Create(comment)
	}
}

func GetComment(doubanId uint64, doubanUid uint64, t uint8) *model.Comment {
	comment := &model.Comment{}
	common.Db.Where("douban_id = ? AND douban_uid = ? AND type = ?", doubanId, doubanUid, t).Find(comment)
	return comment
}

func ListComment(doubanIds *[]uint64, doubanUid uint64, t uint8) *[]model.Comment {
	var comment *[]model.Comment
	common.Db.Where("douban_id IN ? AND douban_uid = ? type = ?", *doubanIds, doubanUid, t).Find(&comment)
	return comment
}
