package controller

import (
	"mouban/internal/consts"
	"mouban/internal/dao"
	"mouban/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RefreshUser godoc
// @Summary      手动刷新用户
// @Tags         admin
// @Produce      json
// @Param        id  query  string  true  "豆瓣用户ID"
// @Success      200  {object}  SuccessOnlyResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /admin/refresh_user [get]
func RefreshUser(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id := util.ParseNumber(idStr)
	if id == 0 {
		BadRequest(ctx, "参数错误")
		return
	}
	schedule := dao.GetSchedule(id, consts.TypeUser.Code)
	if schedule == nil {
		NotFound(ctx, "条目未收录，无法更新")
		return
	}

	user := dao.GetUser(id)
	if user == nil {
		NotFound(ctx, "用户未收录，无法更新")
		return
	}

	dao.RefreshUser(user)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})

	dao.CasScheduleStatus(user.DoubanUid, consts.TypeUser.Code, consts.ScheduleToCrawl.Code, *schedule.Status)

}

// RefreshItem godoc
// @Summary      手动刷新条目
// @Tags         admin
// @Produce      json
// @Param        type  query  string  true  "类型编码"
// @Param        id    query  string  true  "豆瓣条目ID"
// @Success      200  {object}  SuccessOnlyResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      409  {object}  ErrorResponse
// @Router       /admin/refresh_item [get]
func RefreshItem(ctx *gin.Context) {
	typeStr := ctx.Query("type")
	idStr := ctx.Query("id")

	t := uint8(util.ParseNumber(typeStr))
	id := util.ParseNumber(idStr)
	if t == 0 || id == 0 {
		BadRequest(ctx, "参数错误")
		return
	}

	schedule := dao.GetSchedule(id, t)
	if schedule == nil {
		NotFound(ctx, "条目未收录，无法更新")
		return
	}

	if *schedule.Status == consts.ScheduleCrawling.Code || *schedule.Status == consts.ScheduleToCrawl.Code {
		Conflict(ctx, "当前条目正在更新中")
		return
	}

	dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleToCrawl.Code, *schedule.Status)
	logrus.Infoln("refresh item for", t, id)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
