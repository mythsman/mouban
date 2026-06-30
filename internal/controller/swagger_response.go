package controller

import "mouban/internal/model"

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Msg     string `json:"msg" example:"参数错误"`
}

type SuccessOnlyResponse struct {
	Success bool `json:"success" example:"true"`
}

type ResolveUserResult struct {
	Keyword string          `json:"keyword"`
	Users   []ResolveUserVO `json:"users"`
}

type ResolveUserResponse struct {
	Success bool              `json:"success" example:"true"`
	Result  ResolveUserResult `json:"result"`
}

type CheckUserResponse struct {
	Success bool         `json:"success" example:"true"`
	Result  model.UserVO `json:"result"`
}

type BookCommentView struct {
	Item     model.BookVO `json:"item"`
	Rate     uint8        `json:"rate"`
	Label    string       `json:"label"`
	Comment  string       `json:"comment"`
	Action   uint8        `json:"action"`
	MarkDate string       `json:"mark_date"`
}

type MovieCommentView struct {
	Item     model.MovieVO `json:"item"`
	Rate     uint8         `json:"rate"`
	Label    string        `json:"label"`
	Comment  string        `json:"comment"`
	Action   uint8         `json:"action"`
	MarkDate string        `json:"mark_date"`
}

type GameCommentView struct {
	Item     model.GameVO `json:"item"`
	Rate     uint8        `json:"rate"`
	Label    string       `json:"label"`
	Comment  string       `json:"comment"`
	Action   uint8        `json:"action"`
	MarkDate string       `json:"mark_date"`
}

type SongCommentView struct {
	Item     model.SongVO `json:"item"`
	Rate     uint8        `json:"rate"`
	Label    string       `json:"label"`
	Comment  string       `json:"comment"`
	Action   uint8        `json:"action"`
	MarkDate string       `json:"mark_date"`
}

type ListUserBookItemResult struct {
	User    *model.UserVO     `json:"user"`
	Comment []BookCommentView `json:"comment"`
}

type ListUserBookItemResponse struct {
	Success bool                   `json:"success" example:"true"`
	Result  ListUserBookItemResult `json:"result"`
}

type ListUserMovieItemResult struct {
	User    *model.UserVO      `json:"user"`
	Comment []MovieCommentView `json:"comment"`
}

type ListUserMovieItemResponse struct {
	Success bool                    `json:"success" example:"true"`
	Result  ListUserMovieItemResult `json:"result"`
}

type ListUserGameItemResult struct {
	User    *model.UserVO     `json:"user"`
	Comment []GameCommentView `json:"comment"`
}

type ListUserGameItemResponse struct {
	Success bool                   `json:"success" example:"true"`
	Result  ListUserGameItemResult `json:"result"`
}

type ListUserSongItemResult struct {
	User    *model.UserVO     `json:"user"`
	Comment []SongCommentView `json:"comment"`
}

type ListUserSongItemResponse struct {
	Success bool                   `json:"success" example:"true"`
	Result  ListUserSongItemResult `json:"result"`
}

type GuestItemDetailResponse struct {
	Success bool             `json:"success" example:"true"`
	Result  itemDetailResult `json:"result"`
}

type QueueOverviewResponse struct {
	Success bool                `json:"success" example:"true"`
	Result  QueueOverviewResult `json:"result"`
}
