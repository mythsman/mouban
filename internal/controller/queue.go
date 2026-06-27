package controller

import (
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"mouban/internal/consts"
	"mouban/internal/dao"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type QueueTypeOverview struct {
	TypeCode          uint8  `json:"type_code"`
	TypeName          string `json:"type"`
	TypeLabel         string `json:"type_label"`
	ToCrawl           int64  `json:"to_crawl"`
	Crawling          int64  `json:"crawling"`
	CanCrawl          int64  `json:"can_crawl"`
	Unready           int64  `json:"unready"`
	Invalid           int64  `json:"invalid"`
	OldestWaitSeconds int64  `json:"oldest_wait_seconds"`
}

type QueuePoolOverview struct {
	Pool        string  `json:"pool"`
	PoolLabel   string  `json:"pool_label"`
	Concurrency int     `json:"concurrency"`
	Running     int64   `json:"running"`
	Utilization float64 `json:"utilization"`
}

type QueueRunningTaskView struct {
	DoubanID          uint64 `json:"douban_id"`
	TypeCode          uint8  `json:"type_code"`
	TypeName          string `json:"type"`
	TypeLabel         string `json:"type_label"`
	Status            string `json:"status"`
	UpdatedAtUnix     int64  `json:"updated_at_unix"`
	UpdatedAtText     string `json:"updated_at_text"`
	RunningForSeconds int64  `json:"running_for_seconds"`
}

type QueueCompletedTaskView struct {
	DoubanID      uint64 `json:"douban_id"`
	TypeCode      uint8  `json:"type_code"`
	TypeName      string `json:"type"`
	TypeLabel     string `json:"type_label"`
	Status        string `json:"status"`
	Result        string `json:"result"`
	DetailURL     string `json:"detail_url"`
	UpdatedAtUnix int64  `json:"updated_at_unix"`
	UpdatedAtText string `json:"updated_at_text"`
}

type QueueOverviewResult struct {
	Types     []QueueTypeOverview      `json:"types"`
	Pools     []QueuePoolOverview      `json:"pools"`
	Running   []QueueRunningTaskView   `json:"running"`
	Completed []QueueCompletedTaskView `json:"completed"`
}

var queueOverviewCache struct {
	mu       sync.RWMutex
	expireAt time.Time
	data     QueueOverviewResult
}

func QueueOverview(ctx *gin.Context) {
	logAccess(ctx, 0)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  getQueueOverviewCached(20 * time.Second),
	})
}

func getQueueOverviewCached(ttl time.Duration) QueueOverviewResult {
	now := time.Now()
	queueOverviewCache.mu.RLock()
	if now.Before(queueOverviewCache.expireAt) {
		data := queueOverviewCache.data
		queueOverviewCache.mu.RUnlock()
		return data
	}
	queueOverviewCache.mu.RUnlock()

	queueOverviewCache.mu.Lock()
	defer queueOverviewCache.mu.Unlock()
	if time.Now().Before(queueOverviewCache.expireAt) {
		return queueOverviewCache.data
	}
	queueOverviewCache.data = buildQueueOverview()
	queueOverviewCache.expireAt = time.Now().Add(ttl)
	return queueOverviewCache.data
}

func buildQueueOverview() QueueOverviewResult {
	now := time.Now()
	types := []consts.Type{consts.TypeUser, consts.TypeBook, consts.TypeMovie, consts.TypeGame, consts.TypeSong}
	typeMap := map[uint8]*QueueTypeOverview{}
	ordered := make([]QueueTypeOverview, len(types))

	for i, t := range types {
		toCrawl := dao.CountScheduleByTypeAndStatus(t.Code, consts.ScheduleToCrawl.Code)
		crawling := dao.CountScheduleByTypeAndStatus(t.Code, consts.ScheduleCrawling.Code)
		canCrawl := dao.CountScheduleByTypeAndStatus(t.Code, consts.ScheduleCanCrawl.Code)
		unready := dao.CountScheduleByTypeAndResult(t.Code, consts.ScheduleUnready.Code)
		invalid := dao.CountScheduleByTypeAndResult(t.Code, consts.ScheduleInvalid.Code)

		oldestWait := int64(0)
		if oldest := dao.FindOldestUpdatedAtByTypeAndStatus(t.Code, consts.ScheduleToCrawl.Code); oldest != nil {
			oldestWait = int64(now.Sub(*oldest).Seconds())
			if oldestWait < 0 {
				oldestWait = 0
			}
		}

		ordered[i] = QueueTypeOverview{
			TypeCode:          t.Code,
			TypeName:          t.Name,
			TypeLabel:         typeLabel(t.Code),
			ToCrawl:           toCrawl,
			Crawling:          crawling,
			CanCrawl:          canCrawl,
			Unready:           unready,
			Invalid:           invalid,
			OldestWaitSeconds: oldestWait,
		}
		typeMap[t.Code] = &ordered[i]
	}

	runningSchedules := make([]QueueRunningTaskView, 0)
	for _, t := range types {
		rows := dao.ListScheduleByTypeAndStatus(t.Code, consts.ScheduleCrawling.Code, 30)
		for _, row := range rows {
			runningSeconds := int64(now.Sub(row.UpdatedAt).Seconds())
			if runningSeconds < 0 {
				runningSeconds = 0
			}
			statusCode := consts.ScheduleCrawling.Code
			if row.Status != nil {
				statusCode = *row.Status
			}
			runningSchedules = append(runningSchedules, QueueRunningTaskView{
				DoubanID:          row.DoubanId,
				TypeCode:          row.Type,
				TypeName:          consts.ParseType(row.Type).Name,
				TypeLabel:         typeLabel(row.Type),
				Status:            consts.ParseScheduleStatus(statusCode).Name,
				UpdatedAtUnix:     row.UpdatedAt.Unix(),
				UpdatedAtText:     formatTimeCN(row.UpdatedAt),
				RunningForSeconds: runningSeconds,
			})
		}
	}
	sort.Slice(runningSchedules, func(i, j int) bool {
		return runningSchedules[i].UpdatedAtUnix < runningSchedules[j].UpdatedAtUnix
	})
	if len(runningSchedules) > 100 {
		runningSchedules = runningSchedules[:100]
	}

	completedSchedules := make([]QueueCompletedTaskView, 0)
	for _, t := range types {
		rows := dao.ListRecentScheduleByTypeAndStatus(t.Code, consts.ScheduleCrawled.Code, 20)
		for _, row := range rows {
			statusCode := consts.ScheduleCrawled.Code
			if row.Status != nil {
				statusCode = *row.Status
			}
			resultCode := consts.ScheduleUnready.Code
			if row.Result != nil {
				resultCode = *row.Result
			}
			completedSchedules = append(completedSchedules, QueueCompletedTaskView{
				DoubanID:      row.DoubanId,
				TypeCode:      row.Type,
				TypeName:      consts.ParseType(row.Type).Name,
				TypeLabel:     typeLabel(row.Type),
				Status:        consts.ParseScheduleStatus(statusCode).Name,
				Result:        consts.ParseResult(resultCode).Name,
				DetailURL:     buildScheduleDetailURL(row.Type, row.DoubanId),
				UpdatedAtUnix: row.UpdatedAt.Unix(),
				UpdatedAtText: formatTimeCN(row.UpdatedAt),
			})
		}
	}
	sort.Slice(completedSchedules, func(i, j int) bool {
		return completedSchedules[i].UpdatedAtUnix > completedSchedules[j].UpdatedAtUnix
	})
	if len(completedSchedules) > 50 {
		completedSchedules = completedSchedules[:50]
	}

	userConcurrency := viper.GetInt("agent.user.concurrency")
	itemConcurrency := viper.GetInt("agent.item.concurrency")
	userRunning := typeMap[consts.TypeUser.Code].Crawling
	itemRunning := int64(0)
	for _, t := range []consts.Type{consts.TypeBook, consts.TypeMovie, consts.TypeGame, consts.TypeSong} {
		itemRunning += typeMap[t.Code].Crawling
	}

	pools := []QueuePoolOverview{
		{
			Pool:        "user",
			PoolLabel:   "用户队列",
			Concurrency: userConcurrency,
			Running:     userRunning,
			Utilization: calcUtilization(userRunning, userConcurrency),
		},
		{
			Pool:        "item",
			PoolLabel:   "条目队列（book/movie/game/song 共享）",
			Concurrency: itemConcurrency,
			Running:     itemRunning,
			Utilization: calcUtilization(itemRunning, itemConcurrency),
		},
	}

	return QueueOverviewResult{
		Types:     ordered,
		Pools:     pools,
		Running:   runningSchedules,
		Completed: completedSchedules,
	}
}

func calcUtilization(running int64, concurrency int) float64 {
	if concurrency <= 0 {
		return 0
	}
	return float64(running) / float64(concurrency)
}

func typeLabel(code uint8) string {
	switch code {
	case consts.TypeUser.Code:
		return "用户"
	case consts.TypeBook.Code:
		return "图书"
	case consts.TypeMovie.Code:
		return "电影"
	case consts.TypeGame.Code:
		return "游戏"
	case consts.TypeSong.Code:
		return "音乐"
	default:
		return "未知"
	}
}

func buildScheduleDetailURL(t uint8, doubanID uint64) string {
	id := strconv.FormatUint(doubanID, 10)
	switch t {
	case consts.TypeUser.Code:
		return "/user/" + id
	case consts.TypeBook.Code:
		return "/item/book/" + id
	case consts.TypeMovie.Code:
		return "/item/movie/" + id
	case consts.TypeGame.Code:
		return "/item/game/" + id
	case consts.TypeSong.Code:
		return "/item/song/" + id
	default:
		return "#"
	}
}
