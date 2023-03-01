package crawl

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/spf13/viper"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strings"
)

func CommentGame(doubanUid uint64) (*model.User, *[]model.Comment, *[]model.Game, error) {
	var allComments []model.Comment
	var allGames []model.Game
	user := &model.User{}

	comments, games, total := scrollAllGame(doubanUid, consts.ActionDo)
	user.GameDo = total
	allComments = append(allComments, *comments...)
	allGames = append(allGames, *games...)

	comments, games, total = scrollAllGame(doubanUid, consts.ActionWish)
	user.GameWish = total
	allComments = append(allComments, *comments...)
	allGames = append(allGames, *games...)

	comments, games, total = scrollAllGame(doubanUid, consts.ActionCollect)
	user.GameCollect = total
	allComments = append(allComments, *comments...)
	allGames = append(allGames, *games...)

	return user, &allComments, &allGames, nil
}

func scrollAllGame(doubanUid uint64, action consts.Action) (*[]model.Comment, *[]model.Game, uint32) {
	total := uint32(0)
	var allComments []model.Comment
	var allGames []model.Game

	url := ""
	for {
		comments, games, count, next, err := scrollGame(doubanUid, url, action)
		if err != nil {
			panic(err)
		}
		total = count
		url = next
		allComments = append(allComments, *comments...)
		allGames = append(allGames, *games...)
		if next == "" || total >= viper.GetUint32("agent.item.max") {
			break
		}
	}
	return &allComments, &allGames, total
}

func scrollGame(doubanUid uint64, url string, action consts.Action) (*[]model.Comment, *[]model.Game, uint32, string, error) {
	if url == "" {
		url = fmt.Sprintf(consts.GameCommentUrl, doubanUid, action.Name)
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

	list := htmlquery.Find(doc, "//div[@class='common-item']")
	var comments []model.Comment
	var games []model.Game

	for i := range list {
		link := htmlquery.FindOne(list[i], "//div[@class='title']/a")
		href := htmlquery.SelectAttr(link, "href")
		doubanId := util.ParseNumber(href)
		title := htmlquery.InnerText(link)
		thumbnail := htmlquery.SelectAttr(htmlquery.FindOne(list[i], "//div[@class='pic']//img"), "src")

		games = append(games, model.Game{
			DoubanId:  doubanId,
			Title:     strings.TrimSpace(title),
			Thumbnail: thumbnail,
		})

		ratingNumber := uint8(0)
		rating := htmlquery.FindOne(list[i], "//div[@class='rating-info']//span[contains(@class,'rating-star')]")
		if rating != nil {
			ratingNumber = uint8(util.ParseNumber(htmlquery.SelectAttr(rating, "class"))) / 10
		}
		markDate := util.ParseDate(htmlquery.InnerText(htmlquery.FindOne(list[i], "//div[@class='rating-info']//span[@class='date']")))

		shortComment := ""
		shortCommentNode := htmlquery.FindOne(list[i], "/div[@class='content']/div[not(@class)]")
		if shortCommentNode != nil {
			shortComment = util.TrimParagraph(htmlquery.InnerText(shortCommentNode))
		}

		tags := ""
		tagsNode := htmlquery.FindOne(list[i], "//div[@class='rating-info']//span[@class='tags']")
		if tagsNode != nil {
			tags = strings.TrimSpace(strings.Trim(htmlquery.InnerText(tagsNode), "标签:"))
		}

		comment := model.Comment{
			DoubanUid: doubanUid,
			DoubanId:  doubanId,
			Type:      consts.TypeGame.Code,
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
		return &comments, &games, uint32(total), "", err
	} else {
		nextLink := htmlquery.SelectAttr(nextBtn, "href")
		return &comments, &games, uint32(total), fmt.Sprintf("https://www.douban.com/people/%d/games%s", doubanUid, nextLink), err
	}

}
