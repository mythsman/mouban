package crawl

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strings"
)

func CommentMovie(doubanUid uint64) (*model.User, []*model.Comment, error) {
	var allComments []*model.Comment
	user := &model.User{}
	url := ""
	for {
		comments, total, next, err := scrollMovie(doubanUid, url, consts.ActionDo)
		if err != nil {
			return nil, nil, err
		}
		user.MovieDo = total
		url = next
		allComments = append(allComments, comments...)
		if next == "" {
			break
		}
	}

	url = ""
	for {
		comments, total, next, err := scrollMovie(doubanUid, url, consts.ActionWish)
		if err != nil {
			return nil, nil, err
		}
		user.MovieWish = total
		url = next
		allComments = append(allComments, comments...)

		if next == "" {
			break
		}
	}

	url = ""
	for {
		comments, total, next, err := scrollMovie(doubanUid, url, consts.ActionCollect)
		if err != nil {
			return nil, nil, err
		}
		user.MovieCollect = total
		url = next
		allComments = append(allComments, comments...)

		if next == "" {
			break
		}
	}
	return user, allComments, nil
}

func scrollMovie(doubanUid uint64, url string, action consts.Action) ([]*model.Comment, uint32, string, error) {
	if url == "" {
		url = fmt.Sprintf(consts.MovieCommentUrl, doubanUid, action.Name)
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

	list := htmlquery.Find(doc, "//li[@class='subject-item']")
	comments := make([]*model.Comment, len(list))
	for i := range list {
		link := htmlquery.FindOne(list[i], "//div[@class='info']/h2/a")
		href := htmlquery.SelectAttr(link, "href")
		doubanId := util.ParseNumber(href)

		ratingNumber := uint8(0)
		rating := htmlquery.FindOne(list[i], "//div[@class='short-note']//span[contains(@class,'rating')]")
		if rating != nil {
			ratingNumber = uint8(util.ParseNumber(htmlquery.SelectAttr(rating, "class")))
		}
		markDate := util.ParseDate(htmlquery.InnerText(htmlquery.FindOne(list[i], "//div[@class='short-note']//span[@class='date']")))

		shortComment := ""
		shortCommentNode := htmlquery.FindOne(list[i], "//div[@class='short-note']//p[@class='comment']")
		if shortCommentNode != nil {
			shortComment = util.TrimParagraph(htmlquery.InnerText(shortCommentNode))
		}

		tags := ""
		tagsNode := htmlquery.FindOne(list[i], "//div[@class='short-note']//span[@class='tags']")
		if tagsNode != nil {
			tags = strings.TrimSpace(strings.Trim(htmlquery.InnerText(tagsNode), "标签:"))
		}

		comment := &model.Comment{
			DoubanUid: doubanUid,
			DoubanId:  doubanId,
			Type:      consts.TypeMovie,
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
		return comments, uint32(total), "https://movie.douban.com" + nextLink, err
	}

}
