package logic

import (
	"log"
	"mouban/consts"
	"mouban/crawl"
	"mouban/dao"
	"mouban/model"
	"mouban/util"
	"time"
)

func processBook(doubanId uint64) {
	defer func() {
		dao.ChangeScheduleResult(doubanId, consts.TypeBook, consts.ScheduleResultUnready)
	}()
	book, rating, err := crawl.Book(doubanId)

	if err != nil {
		panic(err)
	}
	dao.UpsertBook(book)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeBook, consts.ScheduleResultReady)
}

func processMovie(doubanId uint64) {
	defer func() {
		dao.ChangeScheduleResult(doubanId, consts.TypeMovie, consts.ScheduleResultUnready)
	}()
	movie, rating, err := crawl.Movie(doubanId)

	if err != nil {
		panic(err)
	}
	dao.UpsertMovie(movie)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeMovie, consts.ScheduleResultReady)
}

func processGame(doubanId uint64) {
	defer func() {
		dao.ChangeScheduleResult(doubanId, consts.TypeGame, consts.ScheduleResultUnready)
	}()

	game, rating, err := crawl.Game(doubanId)

	if err != nil {
		panic(err)
	}
	dao.UpsertGame(game)
	dao.UpsertRating(rating)

	dao.ChangeScheduleResult(doubanId, consts.TypeGame, consts.ScheduleResultReady)
}

func processUser(doubanUid uint64) {
	defer func() {
		dao.ChangeScheduleResult(doubanUid, consts.TypeUser, consts.ScheduleResultUnready)
	}()

	//user
	//user, err := crawl.UserOverview(strconv.FormatUint(doubanUid, 10))
	//if err != nil {
	//	panic(err)
	//}
	//dao.UpsertUser(user)

	//game
	_, comment, game, err := crawl.CommentGame(doubanUid)
	if err != nil {
		panic(err)
	}
	for _, g := range *game {
		dao.CreateGameNx(&g)
	}

	for _, c := range *comment {
		dao.UpsertComment(&c)
	}

	////book
	//_, comment, book, err := crawl.CommentBook(doubanUid)
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, b := range *book {
	//	dao.CreateBookNx(&b)
	//}
	//
	//for _, c := range *comment {
	//	dao.UpsertComment(&c)
	//}
	//
	////movie
	//_, comment, movie, err := crawl.CommentMovie(doubanUid)
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, m := range *movie {
	//	dao.CreateMovie(&m)
	//}
	//
	//for _, c := range *comment {
	//	dao.UpsertComment(&c)
	//}
}

func init() {
	ch := make(chan model.Schedule)

	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() {
				if r := recover(); r != nil {
					log.Println(util.GetCurrentGoroutineStack())
				}
			}()

			for {
				schedule := <-ch
				log.Println("agent consume ", schedule)
				switch schedule.Type {
				case consts.TypeBook:
					processBook(schedule.DoubanId)
					break
				case consts.TypeMovie:
					processMovie(schedule.DoubanId)
					break
				case consts.TypeGame:
					processGame(schedule.DoubanId)
					break
				case consts.TypeUser:
					processUser(schedule.DoubanId)
					break
				}
			}
		}(i)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(util.GetCurrentGoroutineStack())
			}
		}()
		for {
			schedule := dao.SearchScheduleByStatus(consts.ScheduleStatusToCrawl)
			if schedule == nil {
				time.Sleep(time.Second * 5)
				log.Println("agent idle")
				continue
			}
			changed := dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawling, consts.ScheduleStatusToCrawl)
			if changed {
				log.Println("agent submit")
				ch <- *schedule
			}
		}
	}()

}
