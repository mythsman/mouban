package dao

import (
	"mouban/common"
	"mouban/model"
)

func UpsertGame(game *model.Game) {
	if common.Db.Where("douban_id = ? ", game.DoubanId).Updates(game).RowsAffected == 0 {
		common.Db.Create(game)
	}
}

func CreateGameNx(game *model.Game) bool {
	return common.Db.Create(game).RowsAffected > 0
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
