package crawl

import (
	"fmt"
	"testing"
)

func TestMovie(t *testing.T) {
	movie, rating, err := Movie(6026035)
	if err != nil {
		return
	}
	fmt.Println(movie)
	fmt.Println(rating)
}

func TestGame(t *testing.T) {
	game, rating, err := Game(26667882)
	if err != nil {
		return
	}
	fmt.Println(game)
	fmt.Println(rating)
}
func TestBook(t *testing.T) {
	book, rating, err := Book(35948443)
	if err != nil {
		return
	}
	fmt.Println(book)
	fmt.Println(rating)

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
	fmt.Println(overview)
}
