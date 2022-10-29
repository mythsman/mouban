package crawl

import (
	"fmt"
	"testing"
)

func TestMovie(t *testing.T) {
	_, _, err := Movie(3908424)
	if err != nil {
		return
	}
}

func TestGame(t *testing.T) {
	_, _, err := Game(10734276)
	if err != nil {
		return
	}
}
func TestBook(t *testing.T) {
	_, _, err := Book(35948443)
	if err != nil {
		return
	}
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
