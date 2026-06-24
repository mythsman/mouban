package crawl

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
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
		panic("download file finally failed")
	}

	mtype, extension := mime(file.Name())

	md5Result := md5sum(file.Name())

	result := upload(file.Name(), md5Result+extension, mtype)

	logrus.Println(result)

	_ = os.Remove(file.Name())
}
