package agent

import (
	"github.com/sirupsen/logrus"
	"mouban/common"
	"mouban/crawl"
	"mouban/dao"
	"mouban/model"
	"mouban/util"
	"time"
)

func work() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorln("flush panic", r, "=>", util.GetCurrentGoroutineStack())
			}
		}()

		rows, _ := common.Db.Model(&model.Game{}).Rows()
		defer rows.Close()

		for rows.Next() {
			var game model.Game
			// ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
			err := common.Db.ScanRows(rows, &game)
			if err != nil {
				logrus.Infoln("game flush to", game.ID)

				newThumb := crawl.Storage(game.Thumbnail)
				if newThumb == game.Thumbnail {
					logrus.Infoln("game flush ignore", game.Thumbnail)
				} else {
					logrus.Infoln("game flush", game.Thumbnail, "->", newThumb)
					game.Thumbnail = newThumb
					dao.UpsertGame(&game)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}
	}()
}
