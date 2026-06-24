package dao

import (
	"mouban/internal/common"
	"mouban/internal/model"
)

func AddAccess(doubanUid uint64, path string, ip string, ua string, referer string) {
	access := &model.Access{
		DoubanUid: doubanUid,
		Path:      path,
		Ip:        ip,
		UserAgent: ua,
		Referer:   referer,
	}

	common.Db.Create(access)
}
