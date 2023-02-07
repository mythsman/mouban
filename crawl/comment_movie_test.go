package crawl

import (
	"github.com/sirupsen/logrus"
	"mouban/consts"
	"mouban/util"
	"testing"
)

func Test_scrollMovie(t *testing.T) {
	comments, movies, total, next, err := scrollMovie(162448367, "https://movie.douban.com/people/mythsman/wish?start=30&sort=time&rating=all&filter=all&mode=grid", consts.ActionCollect)
	if err != nil {
		return
	}
	logrus.Infoln(util.ToJson(*comments))
	logrus.Infoln(util.ToJson(*movies))

	logrus.Infoln(total)
	logrus.Infoln(next)
}
