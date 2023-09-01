package crawl

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strings"
	"time"
)

func CommentBook(user *model.User, forceSyncAfter time.Time) (*[]model.Comment, *[]model.Book, error) {
	var allComments []model.Comment
	var allBooks []model.Book

	if user.BookDo <= viper.GetUint32("agent.item.max") {
		comments, books := scrollAllBook(user.DoubanUid, consts.ActionDo, forceSyncAfter)
		allComments = append(allComments, *comments...)
		allBooks = append(allBooks, *books...)
	}

	if user.BookWish <= viper.GetUint32("agent.item.max") {
		comments, books := scrollAllBook(user.DoubanUid, consts.ActionWish, forceSyncAfter)
		allComments = append(allComments, *comments...)
		allBooks = append(allBooks, *books...)
	}

	if user.BookCollect <= viper.GetUint32("agent.item.max") {
		comments, books := scrollAllBook(user.DoubanUid, consts.ActionCollect, forceSyncAfter)
		allComments = append(allComments, *comments...)
		allBooks = append(allBooks, *books...)
	}

	return &allComments, &allBooks, nil
}

func scrollAllBook(doubanUid uint64, action consts.Action, forceSyncAfter time.Time) (*[]model.Comment, *[]model.Book) {
	total := uint32(0)
	var allComments []model.Comment
	var allBooks []model.Book

	url := ""
	for {
		comments, books, count, next, err := scrollBook(doubanUid, url, action)
		if err != nil {
			panic(err)
		}
		total = count
		url = next
		allComments = append(allComments, *comments...)
		allBooks = append(allBooks, *books...)

		if forceSyncAfter.Unix() > 0 && len(*comments) > 0 {
			if (*comments)[len(*comments)-1].MarkDate.Before(forceSyncAfter) {
				logrus.Infoln("scroll book", action.Name, "for", doubanUid, "end : incr sync end")
				break
			}
		}

		if next == "" {
			logrus.Infoln("scroll book", action.Name, "for", doubanUid, "end : next blank")
			break
		}

		if total >= viper.GetUint32("agent.item.max") {
			logrus.Infoln("scroll book", action.Name, "for", doubanUid, "end : total exceed max")
			break
		}
	}
	return &allComments, &allBooks
}

func scrollBook(doubanUid uint64, url string, action consts.Action) (*[]model.Comment, *[]model.Book, uint32, string, error) {
	if url == "" {
		url = fmt.Sprintf(consts.BookCommentUrl, doubanUid, action.Name)
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
		panic("total is nil for " + url + ", html: " + htmlquery.OutputHTML(doc, true))
	}
	total := util.ParseNumber(htmlquery.InnerText(totalNode))

	list := htmlquery.Find(doc, "//li[@class='subject-item']")
	var comments []model.Comment
	var books []model.Book

	for i := range list {
		link := htmlquery.FindOne(list[i], "//div[@class='info']/h2/a")
		href := htmlquery.SelectAttr(link, "href")
		doubanId := util.ParseNumber(href)
		title := htmlquery.SelectAttr(link, "title")
		thumbnail := htmlquery.SelectAttr(htmlquery.FindOne(list[i], "//div[@class='pic']//img"), "src")

		books = append(books, model.Book{
			DoubanId:  doubanId,
			Title:     strings.TrimSpace(title),
			Thumbnail: thumbnail,
		})

		ratingNumber := uint8(0)
		rating := htmlquery.FindOne(list[i], "//div[@class='short-note']//span[contains(@class,'rating')]")
		if rating != nil {
			ratingNumber = uint8(util.ParseNumber(htmlquery.SelectAttr(rating, "class")))
		}

		markDate := time.Unix(0, 0)
		markDateNode := htmlquery.FindOne(list[i], "//div[@class='short-note']//span[@class='date']")
		if markDateNode != nil {
			markDate = util.ParseDate(htmlquery.InnerText(markDateNode))
		}

		shortComment := ""
		shortCommentNode := htmlquery.FindOne(list[i], "//div[@class='short-note']//p[contains(@class,'comment')]")
		if shortCommentNode != nil {
			shortComment = util.TrimParagraph(htmlquery.InnerText(shortCommentNode))
		}

		tags := ""
		tagsNode := htmlquery.FindOne(list[i], "//div[@class='short-note']//span[@class='tags']")
		if tagsNode != nil {
			tags = strings.TrimSpace(strings.Trim(htmlquery.InnerText(tagsNode), "标签:"))
		}

		comment := model.Comment{
			DoubanUid: doubanUid,
			DoubanId:  doubanId,
			Type:      consts.TypeBook.Code,
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
		return &comments, &books, uint32(total), "", err
	} else {
		nextLink := htmlquery.SelectAttr(nextBtn, "href")
		return &comments, &books, uint32(total), "https://book.douban.com" + nextLink, err
	}

}
