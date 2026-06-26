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

type ItemDetailPageData struct {
	Type      string
	TypeName  string
	ItemID    uint64
	BackURL   string
	DoubanURL string
	Rating    *model.Rating
	Book      *model.Book
	Movie     *model.Movie
	Game      *model.Game
	Song      *model.Song
	Error     string
}

func ItemDetailPage(ctx *gin.Context) {
	typeName := strings.TrimSpace(ctx.Param("type"))
	idStr := strings.TrimSpace(ctx.Param("id"))
	backURL := strings.TrimSpace(ctx.Query("back"))
	if backURL == "" {
		backURL = "/"
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
	logAccess(ctx, 0)
	ctx.HTML(http.StatusOK, "item_detail.tmpl", data)
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
