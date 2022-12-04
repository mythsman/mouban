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

func processBook(doubanId uint64) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, " => ", util.GetCurrentGoroutineStack())
		}
	}()
	book, rating, newUser, err := crawl.Book(doubanId)

	processDiscover(newUser)

	if err != nil {
		dao.ChangeScheduleResult(doubanId, consts.TypeBook, consts.ScheduleResultInvalid)
		panic(err)
	}
	dao.UpsertBook(book)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeBook, consts.ScheduleResultReady)
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
		dao.ChangeScheduleResult(doubanId, consts.TypeMovie, consts.ScheduleResultInvalid)
		panic(err)
	}
	dao.UpsertMovie(movie)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeMovie, consts.ScheduleResultReady)
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
		dao.ChangeScheduleResult(doubanId, consts.TypeGame, consts.ScheduleResultInvalid)
		panic(err)
	}
	dao.UpsertGame(game)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeGame, consts.ScheduleResultReady)
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
				dao.CreateSchedule(id, consts.TypeUser, consts.ScheduleStatusCanCrawl, consts.ScheduleResultUnready)
			}
		}
	}()
}
