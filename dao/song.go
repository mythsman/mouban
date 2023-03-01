package dao

import (
	"mouban/common"
	"mouban/model"
)

func UpsertSong(song *model.Song) {
	data := &model.Song{}
	common.Db.Where("douban_id = ? ", song.DoubanId).Assign(song).FirstOrCreate(data)
}

func CreateSongNx(song *model.Song) bool {
	data := &model.Song{}
	result := common.Db.Where("douban_id = ? ", song.DoubanId).Attrs(song).FirstOrCreate(data)
	return result.RowsAffected > 0
}

func GetSongDetail(doubanId uint64) *model.Song {
	song := &model.Song{}
	common.Db.Where("douban_id = ? ", doubanId).Find(song)
	if song.ID == 0 {
		return nil
	}
	return song
}

func ListSongBrief(doubanIds *[]uint64) *[]model.Song {
	var songs *[]model.Song
	common.Db.Omit("intro", "track_list").Where("douban_id IN ? ", *doubanIds).Find(&songs)
	return songs
}
