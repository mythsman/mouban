package controller

import (
	"net/http"
	"strconv"
	"strings"

	"mouban/internal/consts"
	"mouban/internal/dao"
	"mouban/internal/model"

	"github.com/gin-gonic/gin"
)

type itemRatingView struct {
	Total  uint32  `json:"total"`
	Rating float32 `json:"rating"`
	Star5  float32 `json:"star5"`
	Star4  float32 `json:"star4"`
	Star3  float32 `json:"star3"`
	Star2  float32 `json:"star2"`
	Star1  float32 `json:"star1"`
}

type itemDetailResult struct {
	Type            string          `json:"type"`
	TypeName        string          `json:"type_name"`
	ItemID          uint64          `json:"item_id"`
	DoubanURL       string          `json:"douban_url"`
	CrawledAtText   string          `json:"crawled_at_text"`
	DataUpdatedText string          `json:"data_updated_text"`
	Rating          *itemRatingView `json:"rating"`
	Book            *model.Book     `json:"book,omitempty"`
	Movie           *model.Movie    `json:"movie,omitempty"`
	Game            *model.Game     `json:"game,omitempty"`
	Song            *model.Song     `json:"song,omitempty"`
}

// GuestItemDetail godoc
// @Summary      查询条目详情
// @Tags         guest
// @Produce      json
// @Param        type  query  string  true  "book/movie/game/song"
// @Param        id    query  string  true  "豆瓣条目ID"
// @Success      200  {object}  GuestItemDetailResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /guest/item_detail [get]
func GuestItemDetail(ctx *gin.Context) {
	typeName := strings.TrimSpace(ctx.Query("type"))
	idStr := strings.TrimSpace(ctx.Query("id"))

	itemID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || itemID == 0 {
		BadRequest(ctx, "id 参数错误")
		return
	}

	t, normalizedType, displayName, doubanBase := parseItemType(typeName)
	if t == 0 {
		BadRequest(ctx, "type 参数错误")
		return
	}

	result := itemDetailResult{
		Type:      normalizedType,
		TypeName:  displayName,
		ItemID:    itemID,
		DoubanURL: doubanBase + strconv.FormatUint(itemID, 10) + "/",
	}

	schedule := dao.GetSchedule(itemID, t)
	if schedule != nil {
		result.CrawledAtText = formatTimeCN(schedule.UpdatedAt)
	}
	if result.CrawledAtText == "" {
		result.CrawledAtText = "暂无"
	}

	rating := dao.GetRating(itemID, t)
	if rating != nil {
		result.Rating = &itemRatingView{
			Total:  rating.Total,
			Rating: rating.Rating,
			Star5:  rating.Star5,
			Star4:  rating.Star4,
			Star3:  rating.Star3,
			Star2:  rating.Star2,
			Star1:  rating.Star1,
		}
	}

	switch t {
	case consts.TypeBook.Code:
		book := dao.GetBookDetail(itemID)
		if book == nil {
			NotFound(ctx, "条目不存在或尚未收录")
			return
		}
		result.Book = book
		result.DataUpdatedText = formatTimeCN(book.UpdatedAt)
	case consts.TypeMovie.Code:
		movie := dao.GetMovieDetail(itemID)
		if movie == nil {
			NotFound(ctx, "条目不存在或尚未收录")
			return
		}
		result.Movie = movie
		result.DataUpdatedText = formatTimeCN(movie.UpdatedAt)
	case consts.TypeGame.Code:
		game := dao.GetGameDetail(itemID)
		if game == nil {
			NotFound(ctx, "条目不存在或尚未收录")
			return
		}
		result.Game = game
		result.DataUpdatedText = formatTimeCN(game.UpdatedAt)
	case consts.TypeSong.Code:
		song := dao.GetSongDetail(itemID)
		if song == nil {
			NotFound(ctx, "条目不存在或尚未收录")
			return
		}
		result.Song = song
		result.DataUpdatedText = formatTimeCN(song.UpdatedAt)
	}

	if result.DataUpdatedText == "" {
		result.DataUpdatedText = "暂无"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  result,
	})
}
