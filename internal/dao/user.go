package dao

import (
	"mouban/internal/common"
	"mouban/internal/model"
	"time"

	"github.com/sirupsen/logrus"
)

func CountUser() int64 {
	var count int64
	common.Db.Model(&model.User{}).Count(&count)
	return count
}

func UpsertUser(user *model.User) {
	logrus.WithField("upsert", "user").Infoln("upsert user", user.DoubanUid, user.Name)
	data := &model.User{}
	common.Db.Where("douban_uid = ? ", user.DoubanUid).Assign(user).FirstOrCreate(data)
}

func RefreshUser(user *model.User) {
	logrus.Infoln("refresh user", user.DoubanUid, user.Name)
	common.Db.Model(&model.User{}).
		Where("douban_uid = ? ", user.DoubanUid).
		Updates(model.User{CheckAt: time.Unix(0, 0), SyncAt: time.Unix(0, 0), PublishAt: time.Unix(0, 0)})
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

func ListUserByDomain(domain string) *[]model.User {
	users := []model.User{}
	if domain == "" {
		return &users
	}
	common.Db.Where("domain = ?", domain).Find(&users)
	return &users
}

func ListUserByName(name string) *[]model.User {
	users := []model.User{}
	if name == "" {
		return &users
	}
	common.Db.Where("name = ?", name).Find(&users)
	return &users
}

func ListUserBrief(doubanUids *[]uint64) *[]model.User {
	users := []model.User{}
	if doubanUids == nil || len(*doubanUids) == 0 {
		return &users
	}
	common.Db.Select("douban_uid", "name").Where("douban_uid IN ?", *doubanUids).Find(&users)
	return &users
}
