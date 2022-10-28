package crawl

import "testing"

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
func TestMusic(t *testing.T) {
	_, _, err := Music(5350604)
	if err != nil {
		return
	}
}
