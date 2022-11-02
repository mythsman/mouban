package consts

const RatingNormal = 0
const RatingNotEnough = 1
const RatingNotAllowed = 2

type Action struct {
	Code uint8
	Name string
}

var ActionDo = Action{0, "do"}
var ActionWish = Action{1, "wish"}
var ActionCollect = Action{2, "collect"}

const TypeUser = 0
const TypeBook = 1
const TypeMovie = 2
const TypeGame = 3

const ScheduleStatusToCrawl = 0
const ScheduleStatusCrawling = 1
const ScheduleStatusCrawled = 2

const ScheduleResultUnready = 0
const ScheduleResultReady = 1
const ScheduleResultInvalid = 2
