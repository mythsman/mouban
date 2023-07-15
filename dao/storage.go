package dao

import (
	"github.com/sirupsen/logrus"
	"mouban/common"
	"mouban/model"
)

func UpsertStorage(storage *model.Storage) {
	logrus.Infoln("upsert storage", storage.Source, storage.Target)
	data := &model.Storage{}
	common.Db.Where("source = ? ", storage.Source).Assign(storage).FirstOrCreate(data)
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
