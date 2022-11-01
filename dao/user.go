package dao

import (
	"mouban/common"
	"mouban/model"
)

func UpsertUser(user *model.User) {
	if common.Db.Where("douban_Uid = ? ", user.DoubanUid).Updates(user).RowsAffected == 0 {
		common.Db.Create(user)
	}
}
