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

type Type struct {
	Code uint8
	Name string
}

var TypeUser = Type{0, "user"}
var TypeBook = Type{1, "book"}
var TypeMovie = Type{2, "movie"}
var TypeGame = Type{3, "game"}
var TypeSong = Type{4, "song"}

const ScheduleStatusToCrawl = 0
const ScheduleStatusCrawling = 1
const ScheduleStatusCrawled = 2
const ScheduleStatusCanCrawl = 3

const ScheduleResultUnready = 0
const ScheduleResultReady = 1
const ScheduleResultInvalid = 2
