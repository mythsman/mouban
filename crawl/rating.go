package crawl

import (
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"mouban/consts"
	"mouban/model"
	"mouban/util"
	"strconv"
	"strings"
)

func Rating(interestSelect *html.Node) *model.Rating {
	if interestSelect == nil {
		result := &model.Rating{
			Status: consts.RatingNotAllowed,
		}
		return result
	}

	ratingRaw := htmlquery.InnerText(htmlquery.FindOne(interestSelect, "//strong[@property='v:average']"))
	if len(ratingRaw) == 0 {
		result := &model.Rating{
			Status: consts.RatingNotEnough,
		}
		return result
	}
	rating := util.ParseFloat(ratingRaw)
	totalStr := htmlquery.InnerText(htmlquery.FindOne(interestSelect, "//span[@property='v:votes']"))
	total, err := strconv.ParseUint(totalStr, 10, 32)
	if err != nil {
		return nil
	}
	stars := htmlquery.Find(interestSelect, "//span[@class='rating_per']")

	star5Str := strings.TrimSpace(htmlquery.InnerText(stars[0]))
	star5, err := strconv.ParseFloat(star5Str[0:len(star5Str)-1], 32)
	star4Str := strings.TrimSpace(htmlquery.InnerText(stars[1]))
	star4, err := strconv.ParseFloat(star4Str[0:len(star4Str)-1], 32)
	star3Str := strings.TrimSpace(htmlquery.InnerText(stars[2]))
	star3, err := strconv.ParseFloat(star3Str[0:len(star3Str)-1], 32)
	star2Str := strings.TrimSpace(htmlquery.InnerText(stars[3]))
	star2, err := strconv.ParseFloat(star2Str[0:len(star2Str)-1], 32)
	star1Str := strings.TrimSpace(htmlquery.InnerText(stars[4]))
	star1, err := strconv.ParseFloat(star1Str[0:len(star1Str)-1], 32)

	if err != nil {
		return nil
	}

	result := &model.Rating{
		Total:  uint32(total),
		Rating: rating,
		Star5:  float32(star5),
		Star4:  float32(star4),
		Star3:  float32(star3),
		Star2:  float32(star2),
		Star1:  float32(star1),
		Status: consts.RatingNormal,
	}
	return result
}
