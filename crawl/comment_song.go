package crawl

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/spf13/viper"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strings"
	"time"
)

func CommentSong(user *model.User, forceSyncAfter time.Time) (*[]model.Comment, *[]model.Song, error) {
	var allComments []model.Comment
	var allSongs []model.Song

	if user.SongDo <= viper.GetUint32("agent.item.max") {
		comments, songs := scrollAllSong(user.DoubanUid, consts.ActionDo, forceSyncAfter)
		allComments = append(allComments, *comments...)
		allSongs = append(allSongs, *songs...)
	}

	if user.SongWish <= viper.GetUint32("agent.item.max") {
		comments, songs := scrollAllSong(user.DoubanUid, consts.ActionWish, forceSyncAfter)
		allComments = append(allComments, *comments...)
		allSongs = append(allSongs, *songs...)
	}

	if user.SongCollect <= viper.GetUint32("agent.item.max") {
		comments, songs := scrollAllSong(user.DoubanUid, consts.ActionCollect, forceSyncAfter)
		allComments = append(allComments, *comments...)
		allSongs = append(allSongs, *songs...)
	}
	return &allComments, &allSongs, nil
}

func scrollAllSong(doubanUid uint64, action consts.Action, forceSyncAfter time.Time) (*[]model.Comment, *[]model.Song) {
	total := uint32(0)
	var allComments []model.Comment
	var allSongs []model.Song

	url := ""
	for {
		comments, songs, count, next, err := scrollSong(doubanUid, url, action)
		if err != nil {
			panic(err)
		}
		total = count
		url = next
		allComments = append(allComments, *comments...)
		allSongs = append(allSongs, *songs...)

		if forceSyncAfter.Unix() > 0 && len(*comments) > 0 {
			if (*comments)[len(*comments)-1].MarkDate.Before(forceSyncAfter) {
				break
			}
		}

		if next == "" || total >= viper.GetUint32("agent.item.max") {
			break
		}
	}
	return &allComments, &allSongs
}

func scrollSong(doubanUid uint64, url string, action consts.Action) (*[]model.Comment, *[]model.Song, uint32, string, error) {
	if url == "" {
		url = fmt.Sprintf(consts.SongCommentUrl, doubanUid, action.Name)
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

	list := htmlquery.Find(doc, "//div[@class='item']")
	var comments []model.Comment
	var songs []model.Song

	for i := range list {
		link := htmlquery.FindOne(list[i], "//div[@class='info']//li[@class='title']/a")
		href := htmlquery.SelectAttr(link, "href")
		doubanId := util.ParseNumber(href)
		title := htmlquery.InnerText(htmlquery.FindOne(link, "//em"))
		thumbnail := htmlquery.SelectAttr(htmlquery.FindOne(list[i], "//div[@class='pic']//img"), "src")

		songs = append(songs, model.Song{
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
			Type:      consts.TypeSong.Code,
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
		return &comments, &songs, uint32(total), "", err
	} else {
		nextLink := htmlquery.SelectAttr(nextBtn, "href")
		if !strings.Contains(nextLink, "http") {
			nextLink = "https://music.douban.com" + nextLink
		}
		return &comments, &songs, uint32(total), nextLink, err
	}

}
