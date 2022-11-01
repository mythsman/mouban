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
	user1 := GetUser("1323")
	user2 := GetUser("domain")
	fmt.Println(util.ToJson(user1))
	fmt.Println(util.ToJson(user2))
}
