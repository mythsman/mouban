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
	ItemBaseURL string
	BackURL     string
	Comments    []model.CommentVO
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

type UserExplorerPageData struct {
	Query      string
	SelectedID string
	Candidates []UserCandidateView
	User       *UserProfileView
	Sections   []UserTypeSection
	Error      string
}

func UserExplorerPage(ctx *gin.Context) {
	q := strings.TrimSpace(ctx.Query("q"))
	id := strings.TrimSpace(ctx.Query("id"))

	data := UserExplorerPageData{
		Query:      q,
		SelectedID: id,
		Candidates: []UserCandidateView{},
		Sections:   []UserTypeSection{},
	}

	loggedUserID := uint64(0)

	if id == "" && q != "" {
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

	if id != "" {
		doubanUid, err := strconv.ParseUint(id, 10, 64)
		if err != nil || doubanUid == 0 {
			data.Error = "id 参数错误"
		} else {
			loggedUserID = doubanUid
			user := dao.GetUser(doubanUid)
			if user == nil {
				data.Error = "用户不存在或尚未收录"
			} else {
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
				backURL := buildUserPageURL(q, vo.ID)
				data.Sections = []UserTypeSection{
					buildUserTypeSection(doubanUid, consts.TypeBook.Code, "book", "图书", "想读", "在读", "读过", backURL),
					buildUserTypeSection(doubanUid, consts.TypeMovie.Code, "movie", "电影", "想看", "在看", "看过", backURL),
					buildUserTypeSection(doubanUid, consts.TypeGame.Code, "game", "游戏", "想玩", "在玩", "玩过", backURL),
					buildUserTypeSection(doubanUid, consts.TypeSong.Code, "song", "音乐", "想听", "在听", "听过", backURL),
				}
			}
		}
	}

	logAccess(ctx, loggedUserID)
	ctx.HTML(http.StatusOK, "user_explorer.tmpl", data)
}

func buildUserTypeSection(doubanUid uint64, t uint8, key string, name string, wishLabel string, doLabel string, collectLabel string, backURL string) UserTypeSection {
	wish := buildUserCommentsVO(doubanUid, t, consts.ActionWish.Code, 0)
	doItems := buildUserCommentsVO(doubanUid, t, consts.ActionDo.Code, 0)
	collect := buildUserCommentsVO(doubanUid, t, consts.ActionCollect.Code, 0)

	itemBaseURL := itemBaseURLByType(t)
	return UserTypeSection{
		Key:          key,
		Name:         name,
		WishLabel:    wishLabel,
		DoLabel:      doLabel,
		CollectLabel: collectLabel,
		WishTable:    CommentTableView{ItemBaseURL: itemBaseURL, BackURL: backURL, Comments: wish},
		DoTable:      CommentTableView{ItemBaseURL: itemBaseURL, BackURL: backURL, Comments: doItems},
		CollectTable: CommentTableView{ItemBaseURL: itemBaseURL, BackURL: backURL, Comments: collect},
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

func buildUserPageURL(q string, id uint64) string {
	pageURL := "/?id=" + strconv.FormatUint(id, 10)
	if strings.TrimSpace(q) != "" {
		pageURL += "&q=" + url.QueryEscape(q)
	}
	return pageURL
}

func formatUnixCN(ts int64) string {
	if ts <= 0 {
		return "暂无"
	}
	return time.Unix(ts, 0).In(time.Local).Format("2006-01-02 15:04:05")
}
