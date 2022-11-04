package crawl

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strings"
	"time"
)

func CommentMovie(doubanUid uint64) (*model.User, *[]model.Comment, *[]model.Movie, error) {
	var allComments []model.Comment
	var allMovies []model.Movie

	user := &model.User{}

	comments, movies, total := scrollAllMovie(doubanUid, consts.ActionDo)
	user.MovieDo = total
	allComments = append(allComments, *comments...)
	allMovies = append(allMovies, *movies...)

	comments, movies, total = scrollAllMovie(doubanUid, consts.ActionWish)
	user.MovieWish = total
	allComments = append(allComments, *comments...)
	allMovies = append(allMovies, *movies...)

	comments, movies, total = scrollAllMovie(doubanUid, consts.ActionCollect)
	user.MovieCollect = total
	allComments = append(allComments, *comments...)
	allMovies = append(allMovies, *movies...)

	return user, &allComments, &allMovies, nil
}

func scrollAllMovie(doubanUid uint64, action consts.Action) (*[]model.Comment, *[]model.Movie, uint32) {
	total := uint32(0)
	var allComments []model.Comment
	var allMovies []model.Movie

	url := ""
	for {
		comments, movies, count, next, err := scrollMovie(doubanUid, url, action)
		if err != nil {
			panic(err)
		}
		total = count
		url = next
		allComments = append(allComments, *comments...)
		allMovies = append(allMovies, *movies...)
		if next == "" {
			break
		}
	}
	return &allComments, &allMovies, total
}

func scrollMovie(doubanUid uint64, url string, action consts.Action) (*[]model.Comment, *[]model.Movie, uint32, string, error) {
	if url == "" {
		url = fmt.Sprintf(consts.MovieCommentUrl, doubanUid, action.Name)
	}
	body, err := Get(url)
	if err != nil {
		panic(err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, nil, 0, "", err
	}

	total := util.ParseNumber(htmlquery.InnerText(htmlquery.FindOne(doc, "//div[@id='db-usr-profile']/div[@class='info']/h1")))

	list := htmlquery.Find(doc, "//div[@class='item']")
	var comments []model.Comment
	var movies []model.Movie

	for i := range list {
		link := htmlquery.FindOne(list[i], "//div[@class='info']//li[@class='title']/a")
		href := htmlquery.SelectAttr(link, "href")
		doubanId := util.ParseNumber(href)
		title := htmlquery.InnerText(htmlquery.FindOne(link, "//em"))
		thumbnail := htmlquery.SelectAttr(htmlquery.FindOne(list[i], "//div[@class='pic']//img"), "src")

		movies = append(movies, model.Movie{
			DoubanId:  doubanId,
			Title:     title,
			Thumbnail: thumbnail,
		})

		ratingNumber := uint8(0)
		rating := htmlquery.FindOne(list[i], "//span[contains(@class,'rating')]")
		if rating != nil {
			ratingNumber = uint8(util.ParseNumber(htmlquery.SelectAttr(rating, "class")))
		}
		date := htmlquery.FindOne(list[i], "//span[@class='date']")

		markDate := time.Unix(0, 0)
		if date != nil {
			markDate = util.ParseDate(htmlquery.InnerText(date))
		}

		shortComment := ""
		shortCommentNode := htmlquery.FindOne(list[i], "//span[@class='comment']")
		if shortCommentNode != nil {
			shortComment = util.TrimParagraph(htmlquery.InnerText(shortCommentNode))
		}

		tags := ""
		tagsNode := htmlquery.FindOne(list[i], "//span[@class='tags']")
		if tagsNode != nil {
			tags = strings.TrimSpace(strings.Trim(htmlquery.InnerText(tagsNode), "标签:"))
		}

		comment := model.Comment{
			DoubanUid: doubanUid,
			DoubanId:  doubanId,
			Type:      consts.TypeMovie,
			Rate:      ratingNumber,
			Label:     tags,
			Comment:   shortComment,
			Action:    action.Code,
			MarkDate:  markDate,
		}
		comments = append(comments, comment)
	}

	nextBtn := htmlquery.FindOne(doc, "//link[@rel='next']")
	if nextBtn == nil {
		return &comments, &movies, uint32(total), "", err
	} else {
		nextLink := htmlquery.SelectAttr(nextBtn, "href")
		return &comments, &movies, uint32(total), "https://movie.douban.com" + nextLink, err
	}

}
