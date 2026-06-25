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

type UserTypeSection struct {
	Key          string
	Name         string
	WishLabel    string
	DoLabel      string
	CollectLabel string
	Wish         []model.CommentVO
	Do           []model.CommentVO
	Collect      []model.CommentVO
	Total        int
}

type UserProfileView struct {
	ID            uint64
	Name          string
	Domain        string
	Thumbnail     string
	PublishAtText string
	SyncAtText    string
	CheckAtText   string
}

type UserExplorerPageData struct {
	Query      string
	SelectedID string
	Candidates []ResolveUserVO
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
		Candidates: []ResolveUserVO{},
		Sections:   []UserTypeSection{},
	}

	loggedUserID := uint64(0)

	if id == "" && q != "" {
		data.Candidates = resolveUsers(q)
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
					PublishAtText: formatUnixCN(vo.PublishAt),
					SyncAtText:    formatUnixCN(vo.SyncAt),
					CheckAtText:   formatUnixCN(vo.CheckAt),
				}
				data.Sections = []UserTypeSection{
					buildUserTypeSection(doubanUid, consts.TypeBook.Code, "book", "图书", "想读", "在读", "读过"),
					buildUserTypeSection(doubanUid, consts.TypeMovie.Code, "movie", "电影", "想看", "在看", "看过"),
					buildUserTypeSection(doubanUid, consts.TypeGame.Code, "game", "游戏", "想玩", "在玩", "玩过"),
					buildUserTypeSection(doubanUid, consts.TypeSong.Code, "song", "音乐", "想听", "在听", "听过"),
				}
			}
		}
	}

	logAccess(ctx, loggedUserID)
	ctx.HTML(http.StatusOK, "user_explorer.tmpl", data)
}

func buildUserTypeSection(doubanUid uint64, t uint8, key string, name string, wishLabel string, doLabel string, collectLabel string) UserTypeSection {
	wish := buildUserCommentsVO(doubanUid, t, consts.ActionWish.Code, 0)
	doItems := buildUserCommentsVO(doubanUid, t, consts.ActionDo.Code, 0)
	collect := buildUserCommentsVO(doubanUid, t, consts.ActionCollect.Code, 0)

	return UserTypeSection{
		Key:          key,
		Name:         name,
		WishLabel:    wishLabel,
		DoLabel:      doLabel,
		CollectLabel: collectLabel,
		Wish:         wish,
		Do:           doItems,
		Collect:      collect,
		Total:        len(wish) + len(doItems) + len(collect),
	}
}

func formatUnixCN(ts int64) string {
	if ts <= 0 {
		return "暂无"
	}
	return time.Unix(ts, 0).In(time.Local).Format("2006-01-02 15:04:05")
}
