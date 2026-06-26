package controller

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"mouban/internal/consts"
	"mouban/internal/dao"
	"mouban/internal/model"

	"github.com/gin-gonic/gin"
)

type ItemDetailPageData struct {
	Type            string
	TypeName        string
	ItemID          uint64
	BackURL         string
	DoubanURL       string
	CrawledAtText   string
	DataUpdatedText string
	Rating          *model.Rating
	Book            *model.Book
	Movie           *model.Movie
	Game            *model.Game
	Song            *model.Song
	Error           string
}

func ItemDetailPage(ctx *gin.Context) {
	typeName := strings.TrimSpace(ctx.Param("type"))
	idStr := strings.TrimSpace(ctx.Param("id"))
	backURL := strings.TrimSpace(ctx.Query("back"))
	if backURL == "" {
		backURL = "/explore/users"
	}

	data := ItemDetailPageData{
		Type:    typeName,
		BackURL: backURL,
	}

	itemID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || itemID == 0 {
		data.Error = "条目ID参数错误"
		ctx.HTML(http.StatusBadRequest, "item_detail.tmpl", data)
		return
	}
	data.ItemID = itemID

	t, normalizedType, displayName, doubanBase := parseItemType(typeName)
	if t == 0 {
		data.Error = "条目类型参数错误"
		ctx.HTML(http.StatusBadRequest, "item_detail.tmpl", data)
		return
	}
	data.Type = normalizedType
	data.TypeName = displayName
	data.DoubanURL = doubanBase + strconv.FormatUint(itemID, 10) + "/"

	switch t {
	case consts.TypeBook.Code:
		data.Book = dao.GetBookDetail(itemID)
	case consts.TypeMovie.Code:
		data.Movie = dao.GetMovieDetail(itemID)
	case consts.TypeGame.Code:
		data.Game = dao.GetGameDetail(itemID)
	case consts.TypeSong.Code:
		data.Song = dao.GetSongDetail(itemID)
	}

	if data.Book == nil && data.Movie == nil && data.Game == nil && data.Song == nil {
		data.Error = "条目不存在或尚未收录"
		ctx.HTML(http.StatusNotFound, "item_detail.tmpl", data)
		return
	}

	data.Rating = dao.GetRating(itemID, t)

	schedule := dao.GetSchedule(itemID, t)
	if schedule != nil {
		data.CrawledAtText = formatTimeCN(schedule.UpdatedAt)
	}
	if data.CrawledAtText == "" {
		data.CrawledAtText = "暂无"
	}

	switch t {
	case consts.TypeBook.Code:
		if data.Book != nil {
			data.DataUpdatedText = formatTimeCN(data.Book.UpdatedAt)
		}
	case consts.TypeMovie.Code:
		if data.Movie != nil {
			data.DataUpdatedText = formatTimeCN(data.Movie.UpdatedAt)
		}
	case consts.TypeGame.Code:
		if data.Game != nil {
			data.DataUpdatedText = formatTimeCN(data.Game.UpdatedAt)
		}
	case consts.TypeSong.Code:
		if data.Song != nil {
			data.DataUpdatedText = formatTimeCN(data.Song.UpdatedAt)
		}
	}
	if data.DataUpdatedText == "" {
		data.DataUpdatedText = "暂无"
	}

	logAccess(ctx, 0)
	ctx.HTML(http.StatusOK, "item_detail.tmpl", data)
}

func formatTimeCN(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(time.Local).Format("2006-01-02 15:04:05")
}

func parseItemType(typeName string) (uint8, string, string, string) {
	switch strings.ToLower(strings.TrimSpace(typeName)) {
	case "book":
		return consts.TypeBook.Code, "book", "图书", "https://book.douban.com/subject/"
	case "movie":
		return consts.TypeMovie.Code, "movie", "电影", "https://movie.douban.com/subject/"
	case "game":
		return consts.TypeGame.Code, "game", "游戏", "https://www.douban.com/game/"
	case "song":
		return consts.TypeSong.Code, "song", "音乐", "https://music.douban.com/subject/"
	default:
		return 0, "", "", ""
	}
}
