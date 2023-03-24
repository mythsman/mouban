package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"mouban/consts"
	"mouban/dao"
	"mouban/model"
	"net/http"
	"strconv"
	"time"
)

func CheckUser(ctx *gin.Context) {
	id := ctx.Query("id")
	doubanUid, err := strconv.ParseUint(id, 10, 64)
	if err != nil || id == "0" {
		BizError(ctx, "用户ID输入错误")
		return
	}
	logAccess(ctx, doubanUid)

	schedule := dao.GetSchedule(doubanUid, consts.TypeUser.Code)

	if schedule == nil {
		dao.CreateScheduleNx(doubanUid, consts.TypeUser.Code, consts.ScheduleToCrawl.Code, consts.ScheduleUnready.Code)
		BizError(ctx, "未录入当前用户，已发起录入，请等待后台数据更新")
		return
	}

	if *schedule.Status == consts.ScheduleCanCrawl.Code {
		dao.CasScheduleStatus(doubanUid, consts.TypeUser.Code, consts.ScheduleToCrawl.Code, consts.ScheduleCanCrawl.Code)
		BizError(ctx, "未录入当前用户，已发起录入，请等待后台数据更新")
		return
	}

	if *schedule.Result == consts.ScheduleUnready.Code {
		BizError(ctx, "当前用户录入中")
		return
	}

	if *schedule.Result == consts.ScheduleInvalid.Code {
		BizError(ctx, "当前用户不存在")
		return
	}

	user := dao.GetUser(doubanUid)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  user.Show(),
	})

	if *schedule.Status == consts.ScheduleCrawled.Code {
		timeLimit, _ := time.ParseDuration("-" + viper.GetString("server.limit"))
		if user.CheckAt.Before(time.Now().Add(timeLimit)) {
			dao.CasScheduleStatus(doubanUid, consts.TypeUser.Code, consts.ScheduleToCrawl.Code, consts.ScheduleCrawled.Code)
		}
	}

}

func ListUserItem(ctx *gin.Context, t uint8) {
	id := ctx.Query("id")
	doubanUid, err := strconv.ParseUint(id, 10, 64)
	if err != nil || id == "0" {
		BizError(ctx, "id 参数错误")
		return
	}
	logAccess(ctx, doubanUid)

	action := ctx.Query("action")
	if action == "" {
		BizError(ctx, "action 参数错误")
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
		BizError(ctx, "action 参数错误")
		return
	}

	offset := 0
	if ctx.Query("offset") != "" {
		offset, _ = strconv.Atoi(ctx.Query("offset"))
	}

	user := dao.GetUser(doubanUid)
	if user == nil {
		BizError(ctx, "用户信息找不到")
		return
	}

	comments := dao.SearchComment(doubanUid, t, actionType.Code, offset)

	var ids []uint64
	for _, c := range *comments {
		ids = append(ids, c.DoubanId)
	}

	var commentsVO []model.CommentVO

	switch t {
	case consts.TypeMovie.Code:
		briefs := dao.ListMovieBrief(&ids)
		briefMap := make(map[uint64]*model.Movie)
		for i, _ := range *briefs {
			briefMap[(*briefs)[i].DoubanId] = &(*briefs)[i]
		}

		for i, _ := range *comments {
			movie := briefMap[(*comments)[i].DoubanId]
			commentsVO = append(commentsVO, *(*comments)[i].Show(movie.Show()))
		}
		break
	case consts.TypeBook.Code:
		briefs := dao.ListBookBrief(&ids)
		briefMap := make(map[uint64]*model.Book)
		for i, _ := range *briefs {
			briefMap[(*briefs)[i].DoubanId] = &(*briefs)[i]
		}

		for i, _ := range *comments {
			book := briefMap[(*comments)[i].DoubanId]
			commentsVO = append(commentsVO, *(*comments)[i].Show(book.Show()))
		}
		break
	case consts.TypeGame.Code:
		briefs := dao.ListGameBrief(&ids)
		briefMap := make(map[uint64]*model.Game)
		for i, _ := range *briefs {
			briefMap[(*briefs)[i].DoubanId] = &(*briefs)[i]
		}

		for i, _ := range *comments {
			game := briefMap[(*comments)[i].DoubanId]
			commentsVO = append(commentsVO, *(*comments)[i].Show(game.Show()))
		}
		break
	case consts.TypeSong.Code:
		briefs := dao.ListSongBrief(&ids)
		briefMap := make(map[uint64]*model.Song)
		for i, _ := range *briefs {
			briefMap[(*briefs)[i].DoubanId] = &(*briefs)[i]
		}

		for i, _ := range *comments {
			song := briefMap[(*comments)[i].DoubanId]
			commentsVO = append(commentsVO, *(*comments)[i].Show(song.Show()))
		}
		break
	}

	if commentsVO == nil {
		commentsVO = []model.CommentVO{}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result": gin.H{
			"user":    user.Show(),
			"comment": commentsVO,
		},
	})
}

func logAccess(ctx *gin.Context, doubanUid uint64) {
	ua := ctx.GetHeader("User-Agent")
	referer := ctx.GetHeader("Referer")
	ip := ctx.ClientIP()

	dao.AddAccess(doubanUid, ctx.FullPath(), ip, ua, referer)
}
