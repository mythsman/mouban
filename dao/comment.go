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
