package crawl

import (
	"encoding/xml"
	"errors"
	"fmt"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/sirupsen/logrus"
)

func UserOverview(doubanUid uint64) (*model.User, error) {
	userPublish, err := UserPublish(doubanUid)
	if err != nil {
		return nil, err
	}

	book, err := bookOverview(doubanUid)
	if err != nil {
		return nil, err
	}

	movie, err := movieOverview(doubanUid)
	if err != nil {
		return nil, err
	}

	game, err := gameOverview(doubanUid)
	if err != nil {
		return nil, err

	}

	song, err := songOverview(doubanUid)
	if err != nil {
		return nil, err

	}

	user := &model.User{
		DoubanUid:    book.DoubanUid,
		Domain:       book.Domain,
		Name:         book.Name,
		Thumbnail:    book.Thumbnail,
		BookWish:     book.BookWish,
		BookDo:       book.BookDo,
		BookCollect:  book.BookCollect,
		GameWish:     game.GameWish,
		GameDo:       game.GameDo,
		GameCollect:  game.GameCollect,
		MovieWish:    movie.MovieWish,
		MovieDo:      movie.MovieDo,
		MovieCollect: movie.MovieCollect,
		SongWish:     song.SongWish,
		SongDo:       song.SongDo,
		SongCollect:  song.SongCollect,
		PublishAt:    userPublish,
		RegisterAt:   book.RegisterAt,
	}

	return user, nil
}

func UserPublish(doubanUid uint64) (time.Time, error) {
	body, code, err := Get(fmt.Sprintf(consts.UserRssUrl, doubanUid), UserLimiter)
	if err != nil {
		panic(err)
	}

	if code == 404 {
		return time.Time{}, errors.New("code is 404")
	}

	if code == 403 {
		panic("code is 403 for user rss " + strconv.FormatUint(doubanUid, 10))
	}

	rss := struct {
		XMLName xml.Name `xml:"rss"`
		Channel struct {
			XMLName xml.Name `xml:"channel"`
			Title   string   `xml:"title"`
			PubDate string   `xml:"pubDate"`
		} `xml:"channel"`
	}{}

	data := []byte(*body)

	err = xml.Unmarshal(data, &rss)
	if err != nil {
		logrus.Infoln("rss parse failed for", doubanUid, *body)
		return time.Unix(0, 0), nil
	}

	if rss.Channel.PubDate == "" && rss.Channel.Title != "" {
		return time.Unix(0, 0), nil
	}

	dateTime, err := time.ParseInLocation(time.RFC1123, rss.Channel.PubDate, time.Local)
	if err != nil {
		logrus.Infoln("parse pubDate failed for", doubanUid, *body)
		return time.Unix(0, 0), nil
	}

	return dateTime, nil
}

func UserId(domain string) uint64 {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("user id panic", domain, r, "=>", util.GetCurrentGoroutineStack())
		}
	}()

	body, code, err := Get(fmt.Sprintf(consts.BookOverviewForDomainUrl, domain), DiscoverLimiter)
	if err != nil {
		panic(err)
	}

	if code == 404 {
		logrus.Infoln("user not found for", domain)
		return 0
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		panic(err)
	}

	avatarNode := htmlquery.FindOne(doc, "//div[@id='db-usr-profile']//img")
	if avatarNode == nil {
		panic("avatar not found for " + domain)
	}
	avatarSrc := htmlquery.SelectAttr(avatarNode, "src")
	return util.ParseDoubanUid(avatarSrc)
}

func bookOverview(doubanUid uint64) (*model.User, error) {
	body, code, err := Get(fmt.Sprintf(consts.BookOverviewUrl, doubanUid), UserLimiter)
	if err != nil {
		panic(err)
	}

	if code == 404 {
		logrus.Infoln("user not found for", doubanUid)
		return nil, errors.New("用户不存在")
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, err
	}

	bannedNode := htmlquery.FindOne(doc, "//div[@class='mn']")
	if bannedNode != nil {
		prompt := htmlquery.InnerText(bannedNode)
		if strings.Contains(prompt, "已被永久停用") {
			return nil, errors.New("account banned")
		}
		if strings.Contains(prompt, "已经主动注销") {
			return nil, errors.New("account canceled")
		}
		if strings.Contains(prompt, "已被永久禁言") {
			return nil, errors.New("account forbidden")
		}
	}

	thumbnail := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//div[contains(@class,'book-user-profile')]//img[@class='avatar']"), "src")
	domain := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//div[@id='db-usr-profile']//div[@class='pic']/a"), "href")
	usernameNode := htmlquery.FindOne(doc, "//div[contains(@class,'book-user-profile')]//div[@class='username']")
	if usernameNode == nil {
		panic("username is nil for " + htmlquery.OutputHTML(doc, true))
	}
	username := htmlquery.InnerText(usernameNode)
	registerAt := htmlquery.InnerText(htmlquery.FindOne(doc, "//div[contains(@class,'book-user-profile')]//div[@class='time-registered']"))
	list := htmlquery.Find(doc, "//div[@id='db-book-mine']//h2")
	do := ""
	wish := ""
	collect := ""
	for _, h := range list {
		txt := htmlquery.InnerText(h)
		if strings.Contains(txt, "读过") {
			collect = txt
		}
		if strings.Contains(txt, "想读") {
			wish = txt
		}
		if strings.Contains(txt, "在读") {
			do = txt
		}
	}

	thumbnail = strings.TrimSpace(thumbnail)
	domain = util.ParseDomain(doubanUid, domain)
	username = strings.TrimSpace(username)
	registerTime := util.ParseDate(registerAt)
	doNum := util.ParseNumber(do)
	wishNum := util.ParseNumber(wish)
	collectNum := util.ParseNumber(collect)

	user := &model.User{
		Thumbnail:   thumbnail,
		Domain:      domain,
		DoubanUid:   doubanUid,
		Name:        username,
		RegisterAt:  registerTime,
		BookDo:      uint32(doNum),
		BookWish:    uint32(wishNum),
		BookCollect: uint32(collectNum),
	}
	return user, err

}

func movieOverview(doubanUid uint64) (*model.User, error) {
	body, _, err := Get(fmt.Sprintf(consts.MovieOverviewUrl, doubanUid), UserLimiter)
	if err != nil {
		panic(err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		panic(err)
	}
	domain := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//div[@id='db-usr-profile']//div[@class='pic']/a"), "href")
	domain = util.ParseDomain(doubanUid, domain)

	list := htmlquery.Find(doc, "//div[@id='db-movie-mine']//h2")
	do := ""
	wish := ""
	collect := ""
	for _, h := range list {
		txt := htmlquery.InnerText(h)
		if strings.Contains(txt, "看过") {
			collect = txt
		}
		if strings.Contains(txt, "想看") {
			wish = txt
		}
		if strings.Contains(txt, "在看") {
			do = txt
		}
	}

	doNum := util.ParseNumber(do)
	wishNum := util.ParseNumber(wish)
	collectNum := util.ParseNumber(collect)

	user := &model.User{
		Domain:       domain,
		MovieDo:      uint32(doNum),
		MovieWish:    uint32(wishNum),
		MovieCollect: uint32(collectNum),
	}
	return user, err

}

func songOverview(doubanUid uint64) (*model.User, error) {
	body, _, err := Get(fmt.Sprintf(consts.SongOverviewUrl, doubanUid), UserLimiter)
	if err != nil {
		panic(err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		panic(err)
	}
	domain := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//div[@id='db-usr-profile']//div[@class='pic']/a"), "href")
	domain = util.ParseDomain(doubanUid, domain)

	list := htmlquery.Find(doc, "//div[@id='db-music-mine']//h2")
	do := ""
	wish := ""
	collect := ""
	for _, h := range list {
		txt := htmlquery.InnerText(h)
		if strings.Contains(txt, "听过") {
			collect = txt
		}
		if strings.Contains(txt, "想听") {
			wish = txt
		}
		if strings.Contains(txt, "在听") {
			do = txt
		}
	}

	doNum := util.ParseNumber(do)
	wishNum := util.ParseNumber(wish)
	collectNum := util.ParseNumber(collect)

	user := &model.User{
		Domain:      domain,
		SongDo:      uint32(doNum),
		SongWish:    uint32(wishNum),
		SongCollect: uint32(collectNum),
	}
	return user, err

}

func gameOverview(doubanUid uint64) (*model.User, error) {
	body, _, err := Get(fmt.Sprintf(consts.GameOverviewUrl, doubanUid), UserLimiter)
	if err != nil {
		panic(err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		panic(err)
	}

	domain := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//div[@id='db-usr-profile']//div[@class='pic']/a"), "href")
	domain = util.ParseDomain(doubanUid, domain)

	list := htmlquery.Find(doc, "//div[@class='tabs']//a")
	do := ""
	wish := ""
	collect := ""
	for _, h := range list {
		txt := htmlquery.InnerText(h)
		if strings.Contains(txt, "玩过") {
			collect = txt
		}
		if strings.Contains(txt, "想玩") {
			wish = txt
		}
		if strings.Contains(txt, "在玩") {
			do = txt
		}
	}

	wishNum := util.ParseNumber(wish)
	doNum := util.ParseNumber(do)
	collectNum := util.ParseNumber(collect)

	user := &model.User{
		Domain:      domain,
		GameDo:      uint32(doNum),
		GameWish:    uint32(wishNum),
		GameCollect: uint32(collectNum),
	}
	return user, err

}
