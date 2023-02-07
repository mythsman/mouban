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
	logrus.Info(util.ToJson(*comments))
	logrus.Info(util.ToJson(*movies))

	logrus.Info(total)
	logrus.Info(next)
}
