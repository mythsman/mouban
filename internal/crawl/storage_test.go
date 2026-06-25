package crawl

import (
	"os"
	"strings"
	"testing"
)

func TestStorage(t *testing.T) {
	var file *os.File
	for i := 0; i < 5; i++ {
		file = download("https://img2.doubanio.com/view/subject/s/public/s34828161.jpg", "https://www.douban.com/")
		if file != nil {
			break
		}
	}
	if file == nil {
		t.Fatalf("download file finally failed")
	}
	defer os.Remove(file.Name())

	mtype, extension := mime(file.Name())
	if mtype == "" || extension == "" {
		t.Fatalf("invalid mime detect result: mtype=%q extension=%q", mtype, extension)
	}
	if !strings.HasPrefix(extension, ".") {
		t.Fatalf("extension should start with dot, got %q", extension)
	}

	md5Result := md5sum(file.Name())
	if len(md5Result) != 32 {
		t.Fatalf("invalid md5 len: %d (%s)", len(md5Result), md5Result)
	}
}
