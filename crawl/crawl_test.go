package crawl

import (
	"fmt"
	"testing"
)

func TestMovie(t *testing.T) {
	movie, rating, err := Movie(26235354)
	if err != nil {
		return
	}
	fmt.Println(movie)
	fmt.Println(rating)
}

func TestGame(t *testing.T) {
	_, _, err := Game(10734276)
	if err != nil {
		return
	}
}
func TestBook(t *testing.T) {
	book, rating, err := Book(25863621)
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
