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
		dao.CreateSchedule(doubanUid, consts.TypeUser.Code, consts.ScheduleStatusToCrawl, consts.ScheduleResultUnready)
		BizError(ctx, "未录入当前用户，已发起录入，请等待后台数据更新")
		return
	}

	if schedule.Result == consts.ScheduleResultUnready {
		BizError(ctx, "当前用户录入中")
		return
	}

	if schedule.Result == consts.ScheduleResultInvalid {
		BizError(ctx, "当前用户不存在")
		return
	}

	user := dao.GetUser(doubanUid)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  user.Show(),
	})

	if schedule.Status == consts.ScheduleStatusCrawled {
		timeLimit, _ := time.ParseDuration("-" + viper.GetString("server.limit"))
		if schedule.UpdatedAt.Before(time.Now().Add(timeLimit)) {
			dao.CasScheduleStatus(doubanUid, consts.TypeUser.Code, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
		}
	}

}

func ListUserItem(ctx *gin.Context, t uint8) {
	id := ctx.Query("id")
	doubanUid, err := strconv.ParseUint(id, 10, 64)
	if err != nil || id == "0" {
		BizError(ctx, "用户ID输入错误")
		return
	}
	logAccess(ctx, doubanUid)

	action := ctx.Query("action")
	if action == "" {
		BizError(ctx, "参数错误")
		return
	}

	offset := 0
	if ctx.Query("offset") != "" {
		offset, _ = strconv.Atoi(ctx.Query("offset"))
	}

	schedule := dao.GetSchedule(doubanUid, consts.TypeUser.Code)

	if schedule == nil {
		BizError(ctx, "当前用户未录入")
		return
	}

	if schedule.Result == consts.ScheduleResultUnready {
		BizError(ctx, "当前用户录入中")
		return
	}

	if schedule.Result == consts.ScheduleResultInvalid {
		BizError(ctx, "当前用户不存在")
		return
	}

	user := dao.GetUser(doubanUid)

	comments := dao.SearchComment(doubanUid, t, parseAction(action), offset)

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

func parseAction(action string) uint8 {
	switch action {
	case consts.ActionWish.Name:
		return consts.ActionWish.Code
	case consts.ActionCollect.Name:
		return consts.ActionCollect.Code
	case consts.ActionDo.Name:
		return consts.ActionDo.Code
	}
	return consts.ActionCollect.Code
}
