package crawl

import (
	"mouban/log"
	"mouban/util"
	"testing"
)

func TestMovie(t *testing.T) {
	movie, rating, newUsers, newItems, err := Movie(6026035)
	if err != nil {
		return
	}
	log.Info(util.ToJson(movie))
	log.Info(util.ToJson(rating))
	log.Info(util.ToJson(newUsers))
	log.Info(util.ToJson(newItems))

}

func TestGame(t *testing.T) {
	game, rating, newUsers, newItems, err := Game(35447696)
	if err != nil {
		return
	}
	log.Info(util.ToJson(game))
	log.Info(util.ToJson(rating))
	log.Info(util.ToJson(newUsers))
	log.Info(util.ToJson(newItems))
}

func TestBook(t *testing.T) {
	book, rating, newUser, newItems, err := Book(35948443)
	if err != nil {
		return
	}
	log.Info(util.ToJson(book))
	log.Info(util.ToJson(rating))
	log.Info(util.ToJson(newUser))
	log.Info(util.ToJson(newItems))
}

func TestSong(t *testing.T) {
	//1748967 too many redirects
	song, rating, newUser, newItems, err := Song(1748967)
	if err != nil {
		log.Info(err)
		return
	}
	log.Info(util.ToJson(song))
	log.Info(util.ToJson(rating))
	log.Info(util.ToJson(newUser))
	log.Info(util.ToJson(newItems))
}

func TestUserPublish(t *testing.T) {
	userPublish, err := UserPublish(227565842)
	if err != nil {
		t.Logf("UserPublish failed %s", err)
		return
	}
	t.Logf("UserPublish is %s", userPublish)
}

func TestUserOverview(t *testing.T) {
	overview, err := UserOverview(43001468)
	if err != nil {
		return

	}
	log.Info(util.ToJson(overview))
}

func TestUserId(t *testing.T) {
	id := UserId("162448367")
	log.Info("UserId is", id)
}
