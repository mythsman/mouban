package crawl

import (
	"mouban/internal/util"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMovie(t *testing.T) {
	movie, rating, newUsers, newItems, err := Movie(6026035)
	if err != nil {
		t.Fatalf("Movie failed: %v", err)
	}
	logrus.Infoln(util.ToJson(movie))
	logrus.Infoln(util.ToJson(rating))
	logrus.Infoln(util.ToJson(newUsers))
	logrus.Infoln(util.ToJson(newItems))

}

func TestGame(t *testing.T) {
	game, rating, newUsers, newItems, err := Game(35447696)
	if err != nil {
		t.Fatalf("Game failed: %v", err)
	}
	logrus.Infoln(util.ToJson(game))
	logrus.Infoln(util.ToJson(rating))
	logrus.Infoln(util.ToJson(newUsers))
	logrus.Infoln(util.ToJson(newItems))
}

func TestBook(t *testing.T) {
	book, rating, newUser, newItems, err := Book(35948443)
	if err != nil {
		t.Fatalf("Book failed: %v", err)
	}
	logrus.Infoln(util.ToJson(book))
	logrus.Infoln(util.ToJson(rating))
	logrus.Infoln(util.ToJson(newUser))
	logrus.Infoln(util.ToJson(newItems))
}

func TestSong(t *testing.T) {
	song, rating, newUser, newItems, err := Song(1748967)
	if err != nil {
		if strings.Contains(err.Error(), "too many redirects") {
			t.Skipf("Song skipped due to douban redirect behavior: %v", err)
		}
		t.Fatalf("Song failed: %v", err)
	}
	logrus.Infoln(util.ToJson(song))
	logrus.Infoln(util.ToJson(rating))
	logrus.Infoln(util.ToJson(newUser))
	logrus.Infoln(util.ToJson(newItems))
}

func TestUserPublish(t *testing.T) {
	userPublish, err := UserPublish(227565842)
	if err != nil {
		t.Fatalf("UserPublish failed: %v", err)
	}
	t.Logf("UserPublish is %s", userPublish)
}

func TestUserOverview(t *testing.T) {
	overview, err := UserOverview(43001468)
	if err != nil {
		t.Fatalf("UserOverview failed: %v", err)
	}
	logrus.Infoln(util.ToJson(overview))
}

func TestUserId(t *testing.T) {
	id := UserId("162448367")
	logrus.Infoln("UserId is", id)
}
