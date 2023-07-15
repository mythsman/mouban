package crawl

import (
	"testing"
)

func TestStorage(t *testing.T) {
	res := Storage("https://img1.doubanio.com/lpic/s29499497.jpg")
	t.Log(res)
}
