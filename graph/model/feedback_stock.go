package model

func (FeedbackStock) IsNode() {}

type FeedbackStock struct {
	ID              string `json:"id"`
	IsViewerStocked bool   `json:"isViewerStocked"`
	RecruitmentID   int    `json:"RecruitmentID"`
}
