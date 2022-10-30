package crawl

import (
	"crypto/md5"
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strings"
)

func UserOverview(id string) (*model.User, error) {
	hash, err := UserHash(id)
	if err != nil {
		return nil, err
	}

	book, err := bookOverview(id)
	if err != nil {
		return nil, err
	}

	movie, err := movieOverview(id)
	if err != nil {
		return nil, err
	}
	game, err := gameOverview(id)
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

func UserHash(id string) (string, error) {
	body, err := Get(fmt.Sprintf(consts.UserRssUrl, id))
	if err != nil {
		return "", err
	}
	data := []byte(*body)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)

	return md5str, nil
}

func bookOverview(id string) (*model.User, error) {
	body, err := Get(fmt.Sprintf(consts.BookOverviewUrl, id))
	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, err

	}
	thumbnail := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//div[contains(@class,'book-user-profile')]//img[@class='avatar']"), "src")
	domain := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//div[@id='db-usr-profile']//div[@class='pic']/a"), "href")
	username := htmlquery.InnerText(htmlquery.FindOne(doc, "//div[contains(@class,'book-user-profile')]//div[@class='username']"))
	registerAt := htmlquery.InnerText(htmlquery.FindOne(doc, "//div[contains(@class,'book-user-profile')]//div[@class='time-registered']"))
	list := htmlquery.Find(doc, "//div[@id='db-book-mine']//span[@class='pl']/a")
	do := htmlquery.InnerText(list[0])
	collect := htmlquery.InnerText(list[1])
	wish := htmlquery.InnerText(list[2])

	thumbnail = strings.TrimSpace(thumbnail)
	domain = util.ParseDomain(domain)
	doubanUid := util.ParseDoubanUid(thumbnail)
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

func movieOverview(id string) (*model.User, error) {
	body, err := Get(fmt.Sprintf(consts.MovieOverviewUrl, id))
	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, err
	}
	domain := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//div[@id='db-usr-profile']//div[@class='pic']/a"), "href")
	domain = util.ParseDomain(domain)

	list := htmlquery.Find(doc, "//div[@id='db-movie-mine']//span[@class='pl']/a")
	do := htmlquery.InnerText(list[0])
	collect := htmlquery.InnerText(list[1])
	wish := htmlquery.InnerText(list[2])

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

func gameOverview(id string) (*model.User, error) {
	body, err := Get(fmt.Sprintf(consts.GameOverviewUrl, id))
	if err != nil {
		return nil, err
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, err
	}

	domain := htmlquery.SelectAttr(htmlquery.FindOne(doc, "//div[@id='db-usr-profile']//div[@class='pic']/a"), "href")
	domain = util.ParseDomain(domain)

	list := htmlquery.Find(doc, "//div[@class='tabs']/a")
	wish := htmlquery.InnerText(list[0])
	do := htmlquery.InnerText(list[1])
	collect := htmlquery.InnerText(list[2])

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
