package agent

import (
	"github.com/spf13/viper"
	"log"
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
	}
}

func processBook(doubanId uint64) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, " => ", util.GetCurrentGoroutineStack())
		}
	}()
	book, rating, newUser, err := crawl.Book(doubanId)

	processDiscover(newUser)

	if err != nil {
		dao.ChangeScheduleResult(doubanId, consts.TypeBook.Code, consts.ScheduleResultInvalid)
		panic(err)
	}
	dao.UpsertBook(book)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeBook.Code, consts.ScheduleResultReady)
}

func processMovie(doubanId uint64) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, " => ", util.GetCurrentGoroutineStack())
		}
	}()
	movie, rating, newUser, err := crawl.Movie(doubanId)

	processDiscover(newUser)

	if err != nil {
		dao.ChangeScheduleResult(doubanId, consts.TypeMovie.Code, consts.ScheduleResultInvalid)
		panic(err)
	}
	dao.UpsertMovie(movie)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeMovie.Code, consts.ScheduleResultReady)
}

func processGame(doubanId uint64) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, " => ", util.GetCurrentGoroutineStack())
		}
	}()

	game, rating, newUser, err := crawl.Game(doubanId)

	processDiscover(newUser)

	if err != nil {
		dao.ChangeScheduleResult(doubanId, consts.TypeGame.Code, consts.ScheduleResultInvalid)
		panic(err)
	}
	dao.UpsertGame(game)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeGame.Code, consts.ScheduleResultReady)
}

func processDiscover(newUsers *[]string) {
	if viper.GetString("agent.discover") != "true" {
		return
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r, " => ", util.GetCurrentGoroutineStack())
			}
		}()
		for _, idOrDomain := range *newUsers {
			id, err := strconv.ParseUint(idOrDomain, 10, 64)
			if err != nil {
				user := dao.GetUserByDomain(idOrDomain)
				if user == nil {
					id = crawl.UserId(idOrDomain)
				}
			}
			if id > 0 {
				dao.CreateSchedule(id, consts.TypeUser.Code, consts.ScheduleStatusCanCrawl, consts.ScheduleResultUnready)
			}
		}
	}()
}

func processUser(doubanUid uint64, forceUpdate bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, " => ", util.GetCurrentGoroutineStack())
		}
	}()

	hash, _ := crawl.UserHash(doubanUid)
	rawUser := dao.GetUser(doubanUid)
	if !forceUpdate && rawUser != nil && rawUser.RssHash == hash {
		log.Println("user ", doubanUid, " not changed")
		return
	}
	//user
	user, err := crawl.UserOverview(doubanUid)
	if err != nil {
		dao.ChangeScheduleResult(doubanUid, consts.TypeUser.Code, consts.ScheduleResultInvalid)
		panic(err)
	}
	dao.UpsertUser(user)

	//game
	if user.GameDo > 0 || user.GameWish > 0 || user.GameCollect > 0 {
		_, comment, game, err := crawl.CommentGame(doubanUid)
		if err != nil {
			panic(err)
		}
		go func() {
			for i, _ := range *game {
				dao.UpsertComment(&(*comment)[i])
				added := dao.CreateGameNx(&(*game)[i])
				if added {
					dao.CreateSchedule((*game)[i].DoubanId, consts.TypeGame.Code, consts.ScheduleStatusToCrawl, consts.ScheduleResultUnready)
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
			for i, _ := range *book {
				added := dao.CreateBookNx(&(*book)[i])
				dao.UpsertComment(&(*comment)[i])
				if added {
					dao.CreateSchedule((*book)[i].DoubanId, consts.TypeBook.Code, consts.ScheduleStatusToCrawl, consts.ScheduleResultUnready)
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
			for i, _ := range *movie {
				dao.UpsertComment(&(*comment)[i])
				added := dao.CreateMovieNx(&(*movie)[i])
				if added {
					dao.CreateSchedule((*movie)[i].DoubanId, consts.TypeMovie.Code, consts.ScheduleStatusToCrawl, consts.ScheduleResultUnready)
				}
			}
		}()

	}

	dao.ChangeScheduleResult(doubanUid, consts.TypeUser.Code, consts.ScheduleResultReady)
}
