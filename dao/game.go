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
