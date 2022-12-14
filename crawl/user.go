package crawl

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"log"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strconv"
	"strings"
	"time"
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
		log.Println("rss parse failed for", doubanUid, body)
		return time.Unix(0, 0), nil
	}

	if rss.Channel.PubDate == "" && rss.Channel.Title != "" {
		return time.Unix(0, 0), nil
	}

	dateTime, err := time.ParseInLocation(time.RFC1123, rss.Channel.PubDate, time.Local)
	if err != nil {
		log.Println("parse pubDate failed for", doubanUid, body)
		return time.Unix(0, 0), nil
	}

	return dateTime, nil
}

func UserId(domain string) uint64 {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, " => ", util.GetCurrentGoroutineStack())
		}
	}()

	body, _, err := Get(fmt.Sprintf(consts.BookOverviewForDomainUrl, domain), DiscoverLimiter)
	if err != nil {
		panic(err)
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
	body, _, err := Get(fmt.Sprintf(consts.BookOverviewUrl, doubanUid), UserLimiter)
	if err != nil {
		panic(err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, err
	}

	bannedNode := htmlquery.FindOne(doc, "//div[@class='mn']")
	if bannedNode != nil {
		prompt := htmlquery.InnerText(bannedNode)
		if strings.Contains(prompt, "???????????????????????????") {
			return nil, errors.New("account banned")
		}
		if strings.Contains(prompt, "?????????????????????????????????") {
			return nil, errors.New("account canceled")
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
		if strings.Contains(txt, "??????") {
			collect = txt
		}
		if strings.Contains(txt, "??????") {
			wish = txt
		}
		if strings.Contains(txt, "??????") {
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
		if strings.Contains(txt, "??????") {
			collect = txt
		}
		if strings.Contains(txt, "??????") {
			wish = txt
		}
		if strings.Contains(txt, "??????") {
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
		if strings.Contains(txt, "??????") {
			collect = txt
		}
		if strings.Contains(txt, "??????") {
			wish = txt
		}
		if strings.Contains(txt, "??????") {
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
		if strings.Contains(txt, "??????") {
			collect = txt
		}
		if strings.Contains(txt, "??????") {
			wish = txt
		}
		if strings.Contains(txt, "??????") {
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
