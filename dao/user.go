package dao

import (
	"mouban/common"
	"mouban/model"
)

func UpsertUser(user *model.User) {
	if common.Db.Where("douban_uid = ? ", user.DoubanUid).Updates(user).RowsAffected == 0 {
		common.Db.Create(user)
	}
}

func GetUser(doubanUid uint64) *model.User {
	if doubanUid == 0 {
		return nil
	}
	user := &model.User{}
	common.Db.Where("douban_uid = ? ", doubanUid).Find(user)
	if user.ID == 0 {
		return nil
	}
	return user
}
