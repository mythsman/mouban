package controller

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"mouban/internal/consts"
	"mouban/internal/dao"
	"mouban/internal/model"

	"github.com/gin-gonic/gin"
)

type CommentTableView struct {
	ItemBaseURL   string
	DoubanBaseURL string
	BackURL       string
	Comments      []model.CommentVO
}

type UserTypeSection struct {
	Key          string
	Name         string
	WishLabel    string
	DoLabel      string
	CollectLabel string
	WishTable    CommentTableView
	DoTable      CommentTableView
	CollectTable CommentTableView
	Total        int
}

type UserCandidateView struct {
	ID           uint64
	Name         string
	Domain       string
	Thumbnail    string
	ProfileURL   string
	BookWish     uint32
	BookDo       uint32
	BookCollect  uint32
	GameWish     uint32
	GameDo       uint32
	GameCollect  uint32
	MovieWish    uint32
	MovieDo      uint32
	MovieCollect uint32
	SongWish     uint32
	SongDo       uint32
	SongCollect  uint32
}

type UserProfileView struct {
	ID            uint64
	Name          string
	Domain        string
	Thumbnail     string
	ProfileURL    string
	PublishAtText string
	SyncAtText    string
	CheckAtText   string
}

type UserSearchPageData struct {
	Query      string
	Candidates []UserCandidateView
	Error      string
}

type UserDetailPageData struct {
	Query    string
	User     *UserProfileView
	Sections []UserTypeSection
	Error    string
}

func UserSearchPage(ctx *gin.Context) {
	q := strings.TrimSpace(ctx.Query("q"))
	data := UserSearchPageData{
		Query:      q,
		Candidates: []UserCandidateView{},
	}

	if q != "" {
		resolved := resolveUsers(q)
		for i := range resolved {
			u := resolved[i]
			data.Candidates = append(data.Candidates, UserCandidateView{
				ID:           u.ID,
				Name:         u.Name,
				Domain:       u.Domain,
				Thumbnail:    u.Thumbnail,
				ProfileURL:   buildUserProfileURL(u.ID, u.Domain),
				BookWish:     u.BookWish,
				BookDo:       u.BookDo,
				BookCollect:  u.BookCollect,
				GameWish:     u.GameWish,
				GameDo:       u.GameDo,
				GameCollect:  u.GameCollect,
				MovieWish:    u.MovieWish,
				MovieDo:      u.MovieDo,
				MovieCollect: u.MovieCollect,
				SongWish:     u.SongWish,
				SongDo:       u.SongDo,
				SongCollect:  u.SongCollect,
			})
		}
	}

	logAccess(ctx, 0)
	ctx.HTML(http.StatusOK, "user_search.tmpl", data)
}

func UserDetailPage(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	q := strings.TrimSpace(ctx.Query("q"))

	doubanUID, err := strconv.ParseUint(id, 10, 64)
	if err != nil || doubanUID == 0 {
		ctx.HTML(http.StatusBadRequest, "user_detail.tmpl", UserDetailPageData{
			Query:    q,
			Error:    "id 参数错误",
			Sections: []UserTypeSection{},
		})
		return
	}

	data := UserDetailPageData{
		Query:    q,
		Sections: []UserTypeSection{},
	}

	user := dao.GetUser(doubanUID)
	if user == nil {
		data.Error = "用户不存在或尚未收录"
		logAccess(ctx, doubanUID)
		ctx.HTML(http.StatusNotFound, "user_detail.tmpl", data)
		return
	}

	vo := user.Show()
	data.User = &UserProfileView{
		ID:            vo.ID,
		Name:          vo.Name,
		Domain:        vo.Domain,
		Thumbnail:     vo.Thumbnail,
		ProfileURL:    buildUserProfileURL(vo.ID, vo.Domain),
		PublishAtText: formatUnixCN(vo.PublishAt),
		SyncAtText:    formatUnixCN(vo.SyncAt),
		CheckAtText:   formatUnixCN(vo.CheckAt),
	}

	backURL := buildUserDetailURL(vo.ID, q)
	data.Sections = []UserTypeSection{
		buildUserTypeSection(doubanUID, consts.TypeBook.Code, "book", "图书", "想读", "在读", "读过", backURL),
		buildUserTypeSection(doubanUID, consts.TypeMovie.Code, "movie", "电影", "想看", "在看", "看过", backURL),
		buildUserTypeSection(doubanUID, consts.TypeGame.Code, "game", "游戏", "想玩", "在玩", "玩过", backURL),
		buildUserTypeSection(doubanUID, consts.TypeSong.Code, "song", "音乐", "想听", "在听", "听过", backURL),
	}

	logAccess(ctx, doubanUID)
	ctx.HTML(http.StatusOK, "user_detail.tmpl", data)
}

func buildUserTypeSection(doubanUid uint64, t uint8, key string, name string, wishLabel string, doLabel string, collectLabel string, backURL string) UserTypeSection {
	wish := buildUserCommentsVO(doubanUid, t, consts.ActionWish.Code, 0)
	doItems := buildUserCommentsVO(doubanUid, t, consts.ActionDo.Code, 0)
	collect := buildUserCommentsVO(doubanUid, t, consts.ActionCollect.Code, 0)

	itemBaseURL := itemBaseURLByType(t)
	doubanBaseURL := doubanItemBaseURLByType(t)
	return UserTypeSection{
		Key:          key,
		Name:         name,
		WishLabel:    wishLabel,
		DoLabel:      doLabel,
		CollectLabel: collectLabel,
		WishTable:    CommentTableView{ItemBaseURL: itemBaseURL, DoubanBaseURL: doubanBaseURL, BackURL: backURL, Comments: wish},
		DoTable:      CommentTableView{ItemBaseURL: itemBaseURL, DoubanBaseURL: doubanBaseURL, BackURL: backURL, Comments: doItems},
		CollectTable: CommentTableView{ItemBaseURL: itemBaseURL, DoubanBaseURL: doubanBaseURL, BackURL: backURL, Comments: collect},
		Total:        len(wish) + len(doItems) + len(collect),
	}
}

func buildUserProfileURL(id uint64, domain string) string {
	target := strings.TrimSpace(domain)
	if target == "" {
		target = strconv.FormatUint(id, 10)
	}
	return "https://www.douban.com/people/" + target + "/"
}

func doubanItemBaseURLByType(t uint8) string {
	switch t {
	case consts.TypeBook.Code:
		return "https://book.douban.com/subject/"
	case consts.TypeMovie.Code:
		return "https://movie.douban.com/subject/"
	case consts.TypeGame.Code:
		return "https://www.douban.com/game/"
	case consts.TypeSong.Code:
		return "https://music.douban.com/subject/"
	default:
		return "https://www.douban.com/"
	}
}

func itemBaseURLByType(t uint8) string {
	switch t {
	case consts.TypeBook.Code:
		return "/item/book/"
	case consts.TypeMovie.Code:
		return "/item/movie/"
	case consts.TypeGame.Code:
		return "/item/game/"
	case consts.TypeSong.Code:
		return "/item/song/"
	default:
		return "/"
	}
}

func buildUserSearchURL(q string) string {
	pageURL := "/explore/users"
	if strings.TrimSpace(q) != "" {
		pageURL += "?q=" + url.QueryEscape(q)
	}
	return pageURL
}

func buildUserDetailURL(id uint64, q string) string {
	pageURL := "/user/" + strconv.FormatUint(id, 10)
	if strings.TrimSpace(q) != "" {
		pageURL += "?q=" + url.QueryEscape(q)
	}
	return pageURL
}

func formatUnixCN(ts int64) string {
	if ts <= 0 {
		return "暂无"
	}
	return time.Unix(ts, 0).In(time.Local).Format("2006-01-02 15:04:05")
}
