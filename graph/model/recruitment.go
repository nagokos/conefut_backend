package model

import "time"

func (Recruitment) IsNode() {}

type Recruitment struct {
	ID            string     `json:"id"`
	DatabaseID    int        `json:"databaseId"`
	Title         string     `json:"title"`
	Detail        *string    `json:"detail"`
	Type          Type       `json:"type"`
	Place         *string    `json:"place"`
	StartAt       *time.Time `json:"startAt"`
	LocationLat   *float64   `json:"locationLat"`
	LocationLng   *float64   `json:"locationLng"`
	Status        Status     `json:"status"`
	ClosingAt     *time.Time `json:"closingAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	CreatedAt     time.Time  `json:"createdAt"`
	PublishedAt   *time.Time `json:"publishedAt"`
	UserID        int        `json:"userId"`
	CompetitionID int        `json:"competitionId"`
	PrefectureID  int        `json:"prefectureId"`
}
