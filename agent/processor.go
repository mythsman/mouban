package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mouban/consts"
	"mouban/crawl"
	"mouban/dao"
	"mouban/util"
	"strconv"
)

func processItem(t uint8, doubanId uint64) {
	switch t {
	case consts.TypeBook.Code:
		processBook(doubanId)
		break
	case consts.TypeMovie.Code:
		processMovie(doubanId)
		break
	case consts.TypeGame.Code:
		processGame(doubanId)
		break
	case consts.TypeSong.Code:
		processSong(doubanId)
		break
	}
}

func processBook(doubanId uint64) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, " => ", util.GetCurrentGoroutineStack())
		}
	}()
	book, rating, newUser, newItems, err := crawl.Book(doubanId)

	processDiscoverUser(newUser)
	processDiscoverItem(newItems, consts.TypeBook)

	if err != nil {
		dao.ChangeScheduleResult(doubanId, consts.TypeBook.Code, consts.ScheduleInvalid.Code)
		panic(err)
	}
	dao.UpsertBook(book)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeBook.Code, consts.ScheduleReady.Code)
}

func processMovie(doubanId uint64) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, " => ", util.GetCurrentGoroutineStack())
		}
	}()
	movie, rating, newUser, newItems, err := crawl.Movie(doubanId)

	processDiscoverUser(newUser)
	processDiscoverItem(newItems, consts.TypeMovie)

	if err != nil {
		dao.ChangeScheduleResult(doubanId, consts.TypeMovie.Code, consts.ScheduleInvalid.Code)
		panic(err)
	}
	dao.UpsertMovie(movie)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeMovie.Code, consts.ScheduleReady.Code)
}

func processGame(doubanId uint64) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, " => ", util.GetCurrentGoroutineStack())
		}
	}()

	game, rating, newUser, newItems, err := crawl.Game(doubanId)

	processDiscoverUser(newUser)
	processDiscoverItem(newItems, consts.TypeGame)

	if err != nil {
		dao.ChangeScheduleResult(doubanId, consts.TypeGame.Code, consts.ScheduleInvalid.Code)
		panic(err)
	}
	dao.UpsertGame(game)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeGame.Code, consts.ScheduleReady.Code)
}

func processSong(doubanId uint64) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, " => ", util.GetCurrentGoroutineStack())
		}
	}()

	song, rating, newUser, newItems, err := crawl.Song(doubanId)

	processDiscoverUser(newUser)
	processDiscoverItem(newItems, consts.TypeSong)

	if err != nil {
		dao.ChangeScheduleResult(doubanId, consts.TypeSong.Code, consts.ScheduleInvalid.Code)
		panic(err)
	}
	dao.UpsertSong(song)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeSong.Code, consts.ScheduleReady.Code)
}

func processDiscoverUser(newUsers *[]string) {
	if newUsers == nil {
		return
	}
	level := viper.GetInt("agent.discover.level")
	if level == 0 {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorln(r, " => ", util.GetCurrentGoroutineStack())
			}
		}()
		totalFound := len(*newUsers)
		newFound := 0
		for _, idOrDomain := range *newUsers {
			id, err := strconv.ParseUint(idOrDomain, 10, 64)
			if err != nil {
				if level == 1 {
					continue
				}
				user := dao.GetUserByDomain(idOrDomain)
				if user == nil {
					id = crawl.UserId(idOrDomain)
				}
			}
			if id > 0 {
				added := dao.CreateScheduleNx(id, consts.TypeUser.Code, consts.ScheduleCanCrawl.Code, consts.ScheduleUnready.Code)
				if added {
					newFound += 1
				}
			}
		}
		if newFound > 0 {
			logrus.Infoln("(", newFound, "/", totalFound, ") users discovered")
		}
	}()
}

func processDiscoverItem(newItems *[]uint64, t consts.Type) {
	if newItems == nil || len(*newItems) == 0 {
		return
	}
	level := viper.GetInt("agent.discover.level")
	if level == 0 {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorln(r, " => ", util.GetCurrentGoroutineStack())
			}
		}()
		totalFound := len(*newItems)
		newFound := 0
		for _, doubanId := range *newItems {
			added := dao.CreateScheduleNx(doubanId, t.Code, consts.ScheduleCanCrawl.Code, consts.ScheduleUnready.Code)
			if added {
				newFound += 1
			}
		}
		if newFound > 0 {
			logrus.Infoln("(", newFound, "/", totalFound, ")", t.Name, "discovered")
		}
	}()
}

func processUser(doubanUid uint64) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln(r, " => ", util.GetCurrentGoroutineStack())
		}
	}()

	userPublish, _ := crawl.UserPublish(doubanUid)
	rawUser := dao.GetUser(doubanUid)
	if rawUser != nil && rawUser.PublishAt.Equal(userPublish) {
		logrus.Infoln("user", doubanUid, "not changed")
		return
	}
	//user
	user, err := crawl.UserOverview(doubanUid)
	if err != nil {
		dao.ChangeScheduleResult(doubanUid, consts.TypeUser.Code, consts.ScheduleInvalid.Code)
		panic(err)
	}

	//game
	if user.GameDo > 0 || user.GameWish > 0 || user.GameCollect > 0 {
		_, comment, game, err := crawl.CommentGame(doubanUid)
		if err != nil {
			panic(err)
		}
		go func() {
			newCommentIds := make(map[uint64]bool)
			for i := range *game {
				newCommentIds[(*game)[i].DoubanId] = true
			}
			oldCommentIds := dao.GetCommentIds(doubanUid, consts.TypeGame.Code)
			for i := range *oldCommentIds {
				id := (*oldCommentIds)[i]
				if !newCommentIds[id] {
					dao.HideComment(doubanUid, consts.TypeGame.Code, id)
				}
			}
			for i := range *game {
				dao.UpsertComment(&(*comment)[i])
				added := dao.CreateGameNx(&(*game)[i])
				if added {
					dao.CreateScheduleNx((*game)[i].DoubanId, consts.TypeGame.Code, consts.ScheduleToCrawl.Code, consts.ScheduleUnready.Code)
				}
			}
		}()

	}

	//book
	if user.BookDo > 0 || user.BookWish > 0 || user.BookCollect > 0 {

		_, comment, book, err := crawl.CommentBook(doubanUid)
		if err != nil {
			panic(err)
		}
		go func() {
			newCommentIds := make(map[uint64]bool)
			for i := range *book {
				newCommentIds[(*book)[i].DoubanId] = true
			}
			oldCommentIds := dao.GetCommentIds(doubanUid, consts.TypeBook.Code)
			for i := range *oldCommentIds {
				id := (*oldCommentIds)[i]
				if !newCommentIds[id] {
					dao.HideComment(doubanUid, consts.TypeBook.Code, id)
				}
			}
			for i := range *book {
				added := dao.CreateBookNx(&(*book)[i])
				dao.UpsertComment(&(*comment)[i])
				if added {
					dao.CreateScheduleNx((*book)[i].DoubanId, consts.TypeBook.Code, consts.ScheduleToCrawl.Code, consts.ScheduleUnready.Code)
				}
			}
		}()

	}

	//movie
	if user.MovieDo > 0 || user.MovieWish > 0 || user.MovieCollect > 0 {

		_, comment, movie, err := crawl.CommentMovie(doubanUid)
		if err != nil {
			panic(err)
		}

		go func() {
			newCommentIds := make(map[uint64]bool)
			for i := range *movie {
				newCommentIds[(*movie)[i].DoubanId] = true
			}
			oldCommentIds := dao.GetCommentIds(doubanUid, consts.TypeMovie.Code)
			for i := range *oldCommentIds {
				id := (*oldCommentIds)[i]
				if !newCommentIds[id] {
					dao.HideComment(doubanUid, consts.TypeMovie.Code, id)
				}
			}
			for i := range *movie {
				dao.UpsertComment(&(*comment)[i])
				added := dao.CreateMovieNx(&(*movie)[i])
				if added {
					dao.CreateScheduleNx((*movie)[i].DoubanId, consts.TypeMovie.Code, consts.ScheduleToCrawl.Code, consts.ScheduleUnready.Code)
				}
			}
		}()

	}

	//song
	if user.SongDo > 0 || user.SongWish > 0 || user.SongCollect > 0 {

		_, comment, song, err := crawl.CommentSong(doubanUid)
		if err != nil {
			panic(err)
		}

		go func() {
			newCommentIds := make(map[uint64]bool)
			for i := range *song {
				newCommentIds[(*song)[i].DoubanId] = true
			}
			oldCommentIds := dao.GetCommentIds(doubanUid, consts.TypeSong.Code)
			for i := range *oldCommentIds {
				id := (*oldCommentIds)[i]
				if !newCommentIds[id] {
					dao.HideComment(doubanUid, consts.TypeSong.Code, id)
				}
			}
			for i := range *song {
				dao.UpsertComment(&(*comment)[i])
				added := dao.CreateSongNx(&(*song)[i])
				if added {
					dao.CreateScheduleNx((*song)[i].DoubanId, consts.TypeSong.Code, consts.ScheduleToCrawl.Code, consts.ScheduleUnready.Code)
				}
			}
		}()

	}

	dao.UpsertUser(user)
	dao.ChangeScheduleResult(doubanUid, consts.TypeUser.Code, consts.ScheduleReady.Code)
}
