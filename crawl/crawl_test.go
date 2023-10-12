package crawl

import (
	"mouban/util"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMovie(t *testing.T) {
	movie, rating, newUsers, newItems, err := Movie(6026035)
	if err != nil {
		return
	}
	logrus.Infoln(util.ToJson(movie))
	logrus.Infoln(util.ToJson(rating))
	logrus.Infoln(util.ToJson(newUsers))
	logrus.Infoln(util.ToJson(newItems))

}

func TestGame(t *testing.T) {
	game, rating, newUsers, newItems, err := Game(35447696)
	if err != nil {
		return
	}
	logrus.Infoln(util.ToJson(game))
	logrus.Infoln(util.ToJson(rating))
	logrus.Infoln(util.ToJson(newUsers))
	logrus.Infoln(util.ToJson(newItems))
}

func TestBook(t *testing.T) {
	book, rating, newUser, newItems, err := Book(35948443)
	if err != nil {
		return
	}
	logrus.Infoln(util.ToJson(book))
	logrus.Infoln(util.ToJson(rating))
	logrus.Infoln(util.ToJson(newUser))
	logrus.Infoln(util.ToJson(newItems))
}

func TestSong(t *testing.T) {
	//1748967 too many redirects
	song, rating, newUser, newItems, err := Song(1748967)
	if err != nil {
		logrus.Infoln(err)
		return
	}
	logrus.Infoln(util.ToJson(song))
	logrus.Infoln(util.ToJson(rating))
	logrus.Infoln(util.ToJson(newUser))
	logrus.Infoln(util.ToJson(newItems))
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
	logrus.Infoln(util.ToJson(overview))
}

func TestUserId(t *testing.T) {
	id := UserId("162448367")
	logrus.Infoln("UserId is", id)
}
