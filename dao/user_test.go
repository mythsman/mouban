package dao

import (
	"fmt"
	"mouban/model"
	"mouban/util"
	"testing"
	"time"
)

func TestUpsertUser(t *testing.T) {
	user := &model.User{
		Thumbnail:   "thumbnail",
		Domain:      "domain",
		DoubanUid:   1323,
		Name:        "username",
		RegisterAt:  time.Now(),
		BookDo:      1,
		BookWish:    2,
		BookCollect: 3,
	}
	UpsertUser(user)
}

func TestGetUser(t *testing.T) {
	user := GetUser(1323)
	fmt.Println(util.ToJson(user))
}
