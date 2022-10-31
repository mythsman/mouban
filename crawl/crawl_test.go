package crawl

import (
	"fmt"
	"mouban/util"
	"testing"
)

func TestMovie(t *testing.T) {
	movie, rating, err := Movie(6026035)
	if err != nil {
		return
	}
	fmt.Println(util.ToJson(movie))
	fmt.Println(util.ToJson(rating))
}

func TestGame(t *testing.T) {
	game, rating, err := Game(26667882)
	if err != nil {
		return
	}
	fmt.Println(util.ToJson(game))
	fmt.Println(util.ToJson(rating))
}
func TestBook(t *testing.T) {
	book, rating, err := Book(35948443)
	if err != nil {
		return
	}
	fmt.Println(util.ToJson(book))
	fmt.Println(util.ToJson(rating))
}

func TestUserHash(t *testing.T) {
	hash, err := UserHash("mythsman")
	if err != nil {
		return
	}
	t.Logf("User hash for is %s", hash)
}

func TestUserOverview(t *testing.T) {
	overview, err := UserOverview("mythsman")
	if err != nil {
		return

	}
	fmt.Println(util.ToJson(overview))
}
