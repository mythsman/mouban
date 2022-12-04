package agent

import (
	"github.com/spf13/viper"
	"log"
	"mouban/consts"
	"mouban/crawl"
	"mouban/dao"
	"mouban/util"
	"strconv"
	"time"
)

func processUser(doubanUid uint64) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, " => ", util.GetCurrentGoroutineStack())
		}
	}()

	hash, _ := crawl.UserHash(doubanUid)
	rawUser := dao.GetUser(doubanUid)
	if rawUser != nil && rawUser.RssHash == hash {
		log.Println("user ", doubanUid, " not changed")
		return
	}
	//user
	user, err := crawl.UserOverview(doubanUid)
	if err != nil {
		dao.ChangeScheduleResult(doubanUid, consts.TypeUser, consts.ScheduleResultInvalid)
		panic(err)
	}
	dao.UpsertUser(user)

	//game
	if user.GameDo > 0 || user.GameWish > 0 || user.GameCollect > 0 {
		_, comment, game, err := crawl.CommentGame(doubanUid)
		if err != nil {
			panic(err)
		}
		for i, _ := range *game {
			i := i
			go func() {
				dao.UpsertComment(&(*comment)[i])
				added := dao.CreateGameNx(&(*game)[i])
				if added {
					dao.CreateSchedule((*game)[i].DoubanId, consts.TypeGame, consts.ScheduleStatusToCrawl, consts.ScheduleResultUnready)
				}
			}()
		}
	}

	//
	if user.BookDo > 0 || user.BookWish > 0 || user.BookCollect > 0 {

		_, comment, book, err := crawl.CommentBook(doubanUid)
		if err != nil {
			panic(err)
		}

		for i, _ := range *book {
			i := i
			go func() {
				added := dao.CreateBookNx(&(*book)[i])
				dao.UpsertComment(&(*comment)[i])
				if added {
					dao.CreateSchedule((*book)[i].DoubanId, consts.TypeBook, consts.ScheduleStatusToCrawl, consts.ScheduleResultUnready)
				}
			}()
		}
	}

	//movie
	if user.MovieDo > 0 || user.MovieWish > 0 || user.MovieCollect > 0 {

		_, comment, movie, err := crawl.CommentMovie(doubanUid)
		if err != nil {
			panic(err)
		}

		for i, _ := range *movie {
			i := i
			go func() {
				dao.UpsertComment(&(*comment)[i])
				added := dao.CreateMovieNx(&(*movie)[i])
				if added {
					dao.CreateSchedule((*movie)[i].DoubanId, consts.TypeMovie, consts.ScheduleStatusToCrawl, consts.ScheduleResultUnready)
				}
			}()
		}
	}

	dao.ChangeScheduleResult(doubanUid, consts.TypeUser, consts.ScheduleResultReady)
}
func runUser() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, "user agent crashed  => ", util.GetCurrentGoroutineStack())
		}
	}()
	schedule := dao.SearchScheduleByStatus(consts.TypeUser, consts.ScheduleStatusToCrawl)
	if schedule == nil {
		time.Sleep(time.Second * 5)
	} else {
		changed := dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawling, consts.ScheduleStatusToCrawl)
		if changed {
			log.Println("start process user " + strconv.FormatUint(schedule.DoubanId, 10))
			processUser(schedule.DoubanId)
		}
	}
}
func init() {
	if viper.GetString("agent.enable") != "true" {
		log.Println("user agent disabled")
		return
	}

	go func() {
		for {
			runUser()
		}
	}()

	log.Println("user agent enabled")
}
