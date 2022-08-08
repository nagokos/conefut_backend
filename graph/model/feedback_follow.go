package model

func (FeedbackFollow) IsNode() {}

type FeedbackFollow struct {
	ID               string `json:"id"`
	IsViewerFollowed bool   `json:"isViewerFollowed"`
	UserID           int    `json:"userId"`
}
