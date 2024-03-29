package crawl

import (
	"fmt"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func CommentMovie(user *model.User, forceSyncAfter time.Time) (*[]model.Comment, *[]model.Movie, error) {
	var allComments []model.Comment
	var allMovies []model.Movie

	if user.MovieDo <= viper.GetUint32("agent.item.max") {
		comments, movies := scrollAllMovie(user.DoubanUid, consts.ActionDo, forceSyncAfter)
		allComments = append(allComments, *comments...)
		allMovies = append(allMovies, *movies...)
	}

	if user.MovieWish <= viper.GetUint32("agent.item.max") {
		comments, movies := scrollAllMovie(user.DoubanUid, consts.ActionWish, forceSyncAfter)
		allComments = append(allComments, *comments...)
		allMovies = append(allMovies, *movies...)
	}

	if user.MovieCollect <= viper.GetUint32("agent.item.max") {
		comments, movies := scrollAllMovie(user.DoubanUid, consts.ActionCollect, forceSyncAfter)
		allComments = append(allComments, *comments...)
		allMovies = append(allMovies, *movies...)
	}

	return &allComments, &allMovies, nil
}

func scrollAllMovie(doubanUid uint64, action consts.Action, forceSyncAfter time.Time) (*[]model.Comment, *[]model.Movie) {
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

		if forceSyncAfter.Unix() > 0 && len(*comments) > 0 {
			if (*comments)[len(*comments)-1].MarkDate.Before(forceSyncAfter) {
				logrus.Infoln("scroll movie", action.Name, "for", doubanUid, "end : incr sync end")
				break
			}
		}

		if next == "" {
			logrus.Infoln("scroll movie", action.Name, "for", doubanUid, "end : next blank")
			break
		}

		if total >= viper.GetUint32("agent.item.max") {
			logrus.Infoln("scroll movie", action.Name, "for", doubanUid, "end : total exceed max")
			break
		}
	}
	return &allComments, &allMovies
}

func scrollMovie(doubanUid uint64, url string, action consts.Action) (*[]model.Comment, *[]model.Movie, uint32, string, error) {
	if url == "" {
		url = fmt.Sprintf(consts.MovieCommentUrl, doubanUid, action.Name)
	}
	body, _, err := Get(url, UserLimiter)
	if err != nil {
		panic(err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, nil, 0, "", err
	}

	totalNode := htmlquery.FindOne(doc, "//div[@id='db-usr-profile']/div[@class='info']/h1")
	if totalNode == nil {
		panic("total is nil for " + url + ", html:" + htmlquery.OutputHTML(doc, true))
	}
	total := util.ParseNumber(htmlquery.InnerText(totalNode))

	list := htmlquery.Find(doc, "//div[@class='item']|//div[@class='item comment-item']")
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
			Title:     strings.TrimSpace(title),
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
		shortCommentNode := htmlquery.FindOne(list[i], "//span[contains(@class,'comment')]")
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
			Type:      consts.TypeMovie.Code,
			Rate:      ratingNumber,
			Label:     tags,
			Comment:   shortComment,
			Action:    &action.Code,
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
