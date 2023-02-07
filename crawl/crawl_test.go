package crawl

import (
	"github.com/sirupsen/logrus"
	"mouban/util"
	"testing"
)

func TestMovie(t *testing.T) {
	movie, rating, newUsers, newItems, err := Movie(6026035)
	if err != nil {
		return
	}
	logrus.Info(util.ToJson(movie))
	logrus.Info(util.ToJson(rating))
	logrus.Info(util.ToJson(newUsers))
	logrus.Info(util.ToJson(newItems))

}

func TestGame(t *testing.T) {
	game, rating, newUsers, newItems, err := Game(35447696)
	if err != nil {
		return
	}
	logrus.Info(util.ToJson(game))
	logrus.Info(util.ToJson(rating))
	logrus.Info(util.ToJson(newUsers))
	logrus.Info(util.ToJson(newItems))
}

func TestBook(t *testing.T) {
	book, rating, newUser, newItems, err := Book(35948443)
	if err != nil {
		return
	}
	logrus.Info(util.ToJson(book))
	logrus.Info(util.ToJson(rating))
	logrus.Info(util.ToJson(newUser))
	logrus.Info(util.ToJson(newItems))
}

func TestSong(t *testing.T) {
	//1748967 too many redirects
	song, rating, newUser, newItems, err := Song(1748967)
	if err != nil {
		logrus.Info(err)
		return
	}
	logrus.Info(util.ToJson(song))
	logrus.Info(util.ToJson(rating))
	logrus.Info(util.ToJson(newUser))
	logrus.Info(util.ToJson(newItems))
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
	logrus.Info(util.ToJson(overview))
}

func TestUserId(t *testing.T) {
	id := UserId("162448367")
	logrus.Info("UserId is", id)
}
