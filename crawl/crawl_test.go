package crawl

import (
	"fmt"
	"mouban/util"
	"testing"
)

func TestMovie(t *testing.T) {
	movie, rating, newUsers, err := Movie(6026035)
	if err != nil {
		return
	}
	fmt.Println(util.ToJson(movie))
	fmt.Println(util.ToJson(rating))
	fmt.Println(util.ToJson(newUsers))
}

func TestGame(t *testing.T) {
	game, rating, newUsers, err := Game(35447696)
	if err != nil {
		return
	}
	fmt.Println(util.ToJson(game))
	fmt.Println(util.ToJson(rating))
	fmt.Println(util.ToJson(newUsers))
}

func TestBook(t *testing.T) {
	book, rating, newUser, err := Book(35948443)
	if err != nil {
		return
	}
	fmt.Println(util.ToJson(book))
	fmt.Println(util.ToJson(rating))
	fmt.Println(util.ToJson(newUser))
}

func TestUserHash(t *testing.T) {
	hash, err := UserHash(162448367)
	if err != nil {
		return
	}
	t.Logf("User hash for is %s", hash)
}

func TestUserOverview(t *testing.T) {
	overview, err := UserOverview(235599547)
	if err != nil {
		return

	}
	fmt.Println(util.ToJson(overview))
}

func TestUserId(t *testing.T) {
	id := UserId("162448367")

	t.Logf("UserId is %d", id)
}
