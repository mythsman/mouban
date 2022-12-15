package dao

import (
	"mouban/common"
	"mouban/model"
)

func UpsertSong(song *model.Song) {
	if common.Db.Where("douban_id = ? ", song.DoubanId).Updates(song).RowsAffected == 0 {
		common.Db.Create(song)
	}
}

func CreateSongNx(song *model.Song) bool {
	return common.Db.Create(song).RowsAffected > 0
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
