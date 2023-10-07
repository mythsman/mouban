package dao

import (
	"mouban/common"
	"mouban/model"

	"github.com/sirupsen/logrus"
)

func UpsertStorage(storage *model.Storage) {
	logrus.Infoln("upsert storage", storage.Source, storage.Target)
	data := &model.Storage{}
	common.Db.Where("source = ? ", storage.Source).Assign(storage).FirstOrCreate(data)
}

func GetStorageByMd5(md5 string) *model.Storage {
	if md5 == "" {
		return nil
	}
	storage := &model.Storage{}
	common.Db.Where("md5 = ? ", md5).Limit(1).Find(storage)
	if storage.ID == 0 {
		return nil
	}
	return storage
}

func GetStorage(source string) *model.Storage {
	if source == "" {
		return nil
	}
	storage := &model.Storage{}
	common.Db.Where("source = ? ", source).Find(storage)
	if storage.ID == 0 {
		return nil
	}
	return storage
}
