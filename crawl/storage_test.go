package crawl

import (
	"testing"
)

func TestStorage(t *testing.T) {
	res := Storage("https://img9.doubanio.com/view/subject/l/public/s27244214.jpg")
	t.Log(res)
}
