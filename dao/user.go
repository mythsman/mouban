package dao

import (
	"mouban/common"
	"mouban/model"
	"mouban/util"
	"strconv"
)

func UpsertUser(user *model.User) {
	if common.Db.Where("douban_Uid = ? ", user.DoubanUid).Updates(user).RowsAffected == 0 {
		common.Db.Create(user)
	}
}

func GetUser(doubanUidOrDomain string) *model.User {
	parsable := util.IsParsable(doubanUidOrDomain)
	user := &model.User{}
	if parsable {
		result, _ := strconv.ParseUint(doubanUidOrDomain, 10, 64)
		common.Db.Where("douban_uid = ? ", result).Find(user)
	} else {
		common.Db.Where("domain = ? ", doubanUidOrDomain).Find(user)
	}
	if user.ID == 0 {
		return nil
	}
	return user
}
