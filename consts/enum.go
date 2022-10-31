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

const TypeBook = 0
const TypeMovie = 1
const TypeGame = 2
