package logic

import (
	"mouban/consts"
	"mouban/dao"
)

func Dispatch(doubanId uint64, t uint8) bool {
	schedule := dao.GetSchedule(doubanId, t)
	triggered := false
	switch schedule.Status {
	case consts.ScheduleStatusCrawled:
		dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusToCrawl, consts.ScheduleStatusCrawled)
		triggered = true
		break
	case consts.ScheduleStatusCrawling:
		break
	case consts.ScheduleStatusToCrawl:
		break
	default:
		break
	}
	return triggered
}
