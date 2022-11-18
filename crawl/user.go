package crawl

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strings"
)

func UserOverview(doubanUid uint64) (*model.User, error) {
	hash, err := UserHash(doubanUid)
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
		RssHash:      hash,
		RegisterAt:   book.RegisterAt,
	}

	return user, nil
}

func UserHash(doubanUid uint64) (string, error) {
	body, code, err := Get(fmt.Sprintf(consts.UserRssUrl, doubanUid))
	if err != nil {
		panic(err)
	}

	if code == 404 {
		return "", errors.New("code is 404")
	}

	data := []byte(*body)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)

	return md5str, nil
}

func bookOverview(doubanUid uint64) (*model.User, error) {
	body, _, err := Get(fmt.Sprintf(consts.BookOverviewUrl, doubanUid))
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
		if strings.Contains(prompt, "此帐号已被永久停用") {
			return nil, errors.New("account banned")
		}
	}

	thumbnail := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//div[contains(@class,'book-user-profile')]//img[@class='avatar']"), "src")
	domain := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//div[@id='db-usr-profile']//div[@class='pic']/a"), "href")
	username := htmlquery.InnerText(htmlquery.FindOne(doc, "//div[contains(@class,'book-user-profile')]//div[@class='username']"))
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
	body, _, err := Get(fmt.Sprintf(consts.MovieOverviewUrl, doubanUid))
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

func gameOverview(doubanUid uint64) (*model.User, error) {
	body, _, err := Get(fmt.Sprintf(consts.GameOverviewUrl, doubanUid))
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
