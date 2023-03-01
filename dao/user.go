package dao

import (
	"github.com/sirupsen/logrus"
	"mouban/common"
	"mouban/model"
)

func UpsertUser(user *model.User) {
	logrus.Infoln("upsert user", user.DoubanUid, user.Name)
	data := &model.User{}
	common.Db.Where("douban_uid = ? ", user.DoubanUid).Assign(user).FirstOrCreate(data)
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

func GetUserByDomain(domain string) *model.User {
	if domain == "" {
		return nil
	}
	user := &model.User{}
	common.Db.Where("domain = ? ", domain).Find(user)
	if user.ID == 0 {
		return nil
	}
	return user
}
