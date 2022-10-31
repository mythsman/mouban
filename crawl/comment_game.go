package crawl

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strings"
)

func CommentGame(doubanUid uint64) (*model.User, []*model.Comment, error) {
	var allComments []*model.Comment
	user := &model.User{}
	url := ""
	for {
		comments, total, next, err := scrollGame(doubanUid, url, consts.ActionDo)
		if err != nil {
			return nil, nil, err
		}
		user.GameDo = total
		url = next
		allComments = append(allComments, comments...)
		if next == "" {
			break
		}
	}

	url = ""
	for {
		comments, total, next, err := scrollGame(doubanUid, url, consts.ActionWish)
		if err != nil {
			return nil, nil, err
		}
		user.GameWish = total
		url = next
		allComments = append(allComments, comments...)

		if next == "" {
			break
		}
	}

	url = ""
	for {
		comments, total, next, err := scrollGame(doubanUid, url, consts.ActionCollect)
		if err != nil {
			return nil, nil, err
		}
		user.GameCollect = total
		url = next
		allComments = append(allComments, comments...)

		if next == "" {
			break
		}
	}
	return user, allComments, nil
}

func scrollGame(doubanUid uint64, url string, action consts.Action) ([]*model.Comment, uint32, string, error) {
	if url == "" {
		url = fmt.Sprintf(consts.GameCommentUrl, doubanUid, action.Name)
	}
	body, err := Get(url)
	if err != nil {
		return nil, 0, "", err
	}

	doc, err := htmlquery.Parse(strings.NewReader(*body))
	if err != nil {
		return nil, 0, "", err
	}

	total := util.ParseNumber(htmlquery.InnerText(htmlquery.FindOne(doc, "//div[@id='db-usr-profile']/div[@class='info']/h1")))

	list := htmlquery.Find(doc, "//div[@class='common-item']")
	comments := make([]*model.Comment, len(list))
	for i := range list {
		link := htmlquery.FindOne(list[i], "//div[@class='title']/a")
		href := htmlquery.SelectAttr(link, "href")
		doubanId := util.ParseNumber(href)

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

		comment := &model.Comment{
			DoubanUid: doubanUid,
			DoubanId:  doubanId,
			Type:      consts.TypeGame,
			Rate:      ratingNumber,
			Label:     tags,
			Comment:   shortComment,
			Action:    action.Code,
			MarkDate:  markDate,
		}
		comments[i] = comment
	}

	nextBtn := htmlquery.FindOne(doc, "//link[@rel='next']")
	if nextBtn == nil {
		return comments, uint32(total), "", err
	} else {
		nextLink := htmlquery.SelectAttr(nextBtn, "href")
		return comments, uint32(total), fmt.Sprintf("https://www.douban.com/people/%d/games%s", doubanUid, nextLink), err
	}

}
