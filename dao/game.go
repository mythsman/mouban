package dao

import (
	"github.com/sirupsen/logrus"
	"mouban/common"
	"mouban/model"
)

func UpsertGame(game *model.Game) {
	logrus.Infoln("upsert game", game.DoubanId, game.Title)
	data := &model.Game{}
	common.Db.Where("douban_id = ? ", game.DoubanId).Assign(game).FirstOrCreate(data)
}

func CreateGameNx(game *model.Game) bool {
	data := &model.Game{}
	inserted := common.Db.Where("douban_id = ? ", game.DoubanId).Attrs(game).FirstOrCreate(data).RowsAffected > 0
	if inserted {
		logrus.Infoln("create game", game.DoubanId, game.Title)
	}
	return inserted
}

func GetGameDetail(doubanId uint64) *model.Game {
	game := &model.Game{}
	common.Db.Where("douban_id = ? ", doubanId).Find(game)
	if game.ID == 0 {
		return nil
	}
	return game
}

func ListGameBrief(doubanIds *[]uint64) *[]model.Game {
	var games *[]model.Game
	common.Db.Omit("intro").Where("douban_id IN ? ", *doubanIds).Find(&games)
	return games
}
