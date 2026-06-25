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

type UserTypeSection struct {
	Key     string
	Name    string
	Wish    []model.CommentVO
	Do      []model.CommentVO
	Collect []model.CommentVO
	Total   int
}

type UserExplorerPageData struct {
	Query      string
	SelectedID string
	Candidates []ResolveUserVO
	User       *model.UserVO
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

	if q != "" {
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
				data.User = user.Show()
				data.Sections = []UserTypeSection{
					buildUserTypeSection(doubanUid, consts.TypeBook.Code, "book", "图书"),
					buildUserTypeSection(doubanUid, consts.TypeMovie.Code, "movie", "电影"),
					buildUserTypeSection(doubanUid, consts.TypeGame.Code, "game", "游戏"),
					buildUserTypeSection(doubanUid, consts.TypeSong.Code, "song", "音乐"),
				}
			}
		}
	}

	logAccess(ctx, loggedUserID)
	ctx.HTML(http.StatusOK, "user_explorer.tmpl", data)
}

func buildUserTypeSection(doubanUid uint64, t uint8, key string, name string) UserTypeSection {
	wish := buildUserCommentsVO(doubanUid, t, consts.ActionWish.Code, 0)
	doItems := buildUserCommentsVO(doubanUid, t, consts.ActionDo.Code, 0)
	collect := buildUserCommentsVO(doubanUid, t, consts.ActionCollect.Code, 0)

	return UserTypeSection{
		Key:     key,
		Name:    name,
		Wish:    wish,
		Do:      doItems,
		Collect: collect,
		Total:   len(wish) + len(doItems) + len(collect),
	}
}
