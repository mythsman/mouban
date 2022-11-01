package dao

import (
	"mouban/common"
	"mouban/model"
)

func UpsertBook(book *model.Book) {
	if common.Db.Where("douban_id = ? ", book.DoubanId).Updates(book).RowsAffected == 0 {
		common.Db.Create(book)
	}
}
