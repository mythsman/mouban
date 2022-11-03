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

func CreateBookNx(book *model.Book) bool {
	return common.Db.Create(book).RowsAffected > 0
}

func GetBookDetail(doubanId uint64) *model.Book {
	book := &model.Book{}
	common.Db.Where("douban_id = ? ", doubanId).Find(book)
	if book.ID == 0 {
		return nil
	}
	return book
}

func ListBookBrief(doubanIds *[]uint64) *[]model.Book {
	var books *[]model.Book
	common.Db.Omit("serial", "isbn", "framing", "page", "intro").Where("douban_id IN ? ", *doubanIds).Find(&books)
	return books
}
