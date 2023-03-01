package dao

import (
	"github.com/sirupsen/logrus"
	"mouban/common"
	"mouban/model"
)

func UpsertBook(book *model.Book) {
	logrus.Infoln("upsert book", book.DoubanId, book.Title)
	data := &model.Book{}
	common.Db.Where("douban_id = ? ", book.DoubanId).Assign(book).FirstOrCreate(data)
}

func CreateBookNx(book *model.Book) bool {
	logrus.Infoln("create book", book.DoubanId, book.Title)
	data := &model.Book{}
	result := common.Db.Where("douban_id = ? ", book.DoubanId).Attrs(book).FirstOrCreate(data)
	return result.RowsAffected > 0
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
