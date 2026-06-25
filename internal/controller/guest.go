package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mouban/internal/consts"
	"mouban/internal/dao"
	"mouban/internal/model"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ResolveUserVO struct {
	model.UserVO
	MatchBy []string `json:"match_by"`
}

type resolveUserCandidate struct {
	user    *model.User
	matchBy map[string]bool
}

func ResolveUser(ctx *gin.Context) {
	q := strings.TrimSpace(ctx.Query("q"))
	if q == "" {
		BadRequest(ctx, "q 参数错误")
		return
	}

	logAccess(ctx, 0)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result": gin.H{
			"keyword": q,
			"users":   resolveUsers(q),
		},
	})
}

func CheckUser(ctx *gin.Context) {
	id := ctx.Query("id")
	doubanUid, err := strconv.ParseUint(id, 10, 64)
	if err != nil || id == "0" {
		BadRequest(ctx, "用户ID输入错误")
		return
	}
	logAccess(ctx, doubanUid)

	schedule := dao.GetSchedule(doubanUid, consts.TypeUser.Code)

	if schedule == nil {
		dao.CreateScheduleNx(doubanUid, consts.TypeUser.Code, consts.ScheduleToCrawl.Code, consts.ScheduleUnready.Code)
		Accepted(ctx, "未录入当前用户，已发起录入，请等待后台数据更新")
		return
	}

	if *schedule.Status == consts.ScheduleCanCrawl.Code {
		dao.CasScheduleStatus(doubanUid, consts.TypeUser.Code, consts.ScheduleToCrawl.Code, consts.ScheduleCanCrawl.Code)
		Accepted(ctx, "未录入当前用户，已发起录入，请等待后台数据更新")
		return
	}

	if *schedule.Result == consts.ScheduleUnready.Code {
		Accepted(ctx, "当前用户录入中")
		return
	}

	if *schedule.Result == consts.ScheduleInvalid.Code {
		NotFound(ctx, "当前用户不存在")
		return
	}

	user := dao.GetUser(doubanUid)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  user.Show(),
	})

	if *schedule.Status == consts.ScheduleCrawled.Code {
		timeLimit, _ := time.ParseDuration("-" + viper.GetString("user.recheck_interval"))
		if user.CheckAt.Before(time.Now().Add(timeLimit)) {
			dao.CasScheduleStatus(doubanUid, consts.TypeUser.Code, consts.ScheduleToCrawl.Code, consts.ScheduleCrawled.Code)
		}
	}

}

func ListUserItem(ctx *gin.Context, t uint8) {
	id := ctx.Query("id")
	doubanUid, err := strconv.ParseUint(id, 10, 64)
	if err != nil || id == "0" {
		BadRequest(ctx, "id 参数错误")
		return
	}
	logAccess(ctx, doubanUid)

	action := ctx.Query("action")
	if action == "" {
		BadRequest(ctx, "action 参数错误")
		return
	}

	var actionType *consts.Action
	switch action {
	case consts.ActionWish.Name:
		actionType = &consts.ActionWish
	case consts.ActionCollect.Name:
		actionType = &consts.ActionCollect
	case consts.ActionDo.Name:
		actionType = &consts.ActionDo
	case consts.ActionHide.Name:
		actionType = &consts.ActionHide
	}
	if actionType == nil {
		BadRequest(ctx, "action 参数错误")
		return
	}

	offset := 0
	if ctx.Query("offset") != "" {
		offset, _ = strconv.Atoi(ctx.Query("offset"))
	}

	user := dao.GetUser(doubanUid)
	if user == nil {
		NotFound(ctx, "用户信息找不到")
		return
	}

	commentsVO := buildUserCommentsVO(doubanUid, t, actionType.Code, offset)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result": gin.H{
			"user":    user.Show(),
			"comment": commentsVO,
		},
	})
}

func resolveUsers(q string) []ResolveUserVO {
	candidates := make(map[uint64]*resolveUserCandidate)
	order := make([]uint64, 0)

	addUser := func(user *model.User, by string) {
		if user == nil {
			return
		}
		if user.DoubanUid == 0 {
			return
		}
		candidate, ok := candidates[user.DoubanUid]
		if !ok {
			candidate = &resolveUserCandidate{user: user, matchBy: map[string]bool{}}
			candidates[user.DoubanUid] = candidate
			order = append(order, user.DoubanUid)
		}
		candidate.matchBy[by] = true
	}

	addUsers := func(users *[]model.User, by string) {
		for i := range *users {
			addUser(&(*users)[i], by)
		}
	}

	if doubanUid, err := strconv.ParseUint(q, 10, 64); err == nil && doubanUid > 0 {
		addUser(dao.GetUser(doubanUid), "id")
	}
	addUsers(dao.ListUserByDomain(q), "domain")
	addUsers(dao.ListUserByName(q), "name")

	result := make([]ResolveUserVO, 0, len(order))
	for _, doubanUid := range order {
		candidate := candidates[doubanUid]
		if candidate == nil || candidate.user == nil {
			continue
		}
		matchBy := make([]string, 0, 3)
		for _, by := range []string{"id", "domain", "name"} {
			if candidate.matchBy[by] {
				matchBy = append(matchBy, by)
			}
		}
		result = append(result, ResolveUserVO{UserVO: *candidate.user.Show(), MatchBy: matchBy})
	}
	return result
}

func buildUserCommentsVO(doubanUid uint64, t uint8, actionCode uint8, offset int) []model.CommentVO {
	comments := dao.SearchComment(doubanUid, t, actionCode, offset)

	ids := make([]uint64, 0, len(*comments))
	for _, c := range *comments {
		ids = append(ids, c.DoubanId)
	}

	commentsVO := make([]model.CommentVO, 0, len(*comments))
	if len(ids) == 0 {
		return commentsVO
	}

	switch t {
	case consts.TypeMovie.Code:
		briefs := dao.ListMovieBrief(&ids)
		briefMap := make(map[uint64]*model.Movie, len(*briefs))
		for i := range *briefs {
			briefMap[(*briefs)[i].DoubanId] = &(*briefs)[i]
		}
		for i := range *comments {
			movie := briefMap[(*comments)[i].DoubanId]
			if movie != nil {
				commentsVO = append(commentsVO, *(*comments)[i].Show(movie.Show()))
			}
		}
	case consts.TypeBook.Code:
		briefs := dao.ListBookBrief(&ids)
		briefMap := make(map[uint64]*model.Book, len(*briefs))
		for i := range *briefs {
			briefMap[(*briefs)[i].DoubanId] = &(*briefs)[i]
		}
		for i := range *comments {
			book := briefMap[(*comments)[i].DoubanId]
			if book != nil {
				commentsVO = append(commentsVO, *(*comments)[i].Show(book.Show()))
			}
		}
	case consts.TypeGame.Code:
		briefs := dao.ListGameBrief(&ids)
		briefMap := make(map[uint64]*model.Game, len(*briefs))
		for i := range *briefs {
			briefMap[(*briefs)[i].DoubanId] = &(*briefs)[i]
		}
		for i := range *comments {
			game := briefMap[(*comments)[i].DoubanId]
			if game != nil {
				commentsVO = append(commentsVO, *(*comments)[i].Show(game.Show()))
			}
		}
	case consts.TypeSong.Code:
		briefs := dao.ListSongBrief(&ids)
		briefMap := make(map[uint64]*model.Song, len(*briefs))
		for i := range *briefs {
			briefMap[(*briefs)[i].DoubanId] = &(*briefs)[i]
		}
		for i := range *comments {
			song := briefMap[(*comments)[i].DoubanId]
			if song != nil {
				commentsVO = append(commentsVO, *(*comments)[i].Show(song.Show()))
			}
		}
	}

	return commentsVO
}

func logAccess(ctx *gin.Context, doubanUid uint64) {
	ua := ctx.GetHeader("User-Agent")
	referer := ctx.GetHeader("Referer")
	ip := ctx.ClientIP()

	dao.AddAccess(doubanUid, ctx.FullPath(), ip, ua, referer)
}
