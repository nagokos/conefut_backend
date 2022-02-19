package recruitment

import (
	"context"
	"fmt"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/ent/recruitment"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

type Recruitment struct {
	Title         string
	Type          model.Type
	Level         model.Level
	Place         *string
	StartAt       *time.Time
	Content       *string
	LocationLat   *float64
	LocationLng   *float64
	IsPublished   bool
	Capacity      *int
	ClosingAt     *time.Time
	CompetitionID *string
	PrefectureID  *string
}

func (r Recruitment) CreateRecruitmentValidate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.Title,
			validation.Required.Error("タイトルを入力してください"),
			validation.RuneLength(1, 60).Error("タイトルは60文字以内で入力してください"),
		),
		// validation.Field(
		// 	&r.Type,
		// 	validation.Required.Error("募集タイプを選択してください"),
		// ),
		// validation.Field(
		// 	&r.Content,
		// 	validation.Required.Error("募集の詳細を入力してください"),
		// 	validation.RuneLength(1, 10000).Error("募集の詳細は10000文字以内で入力してください"),
		// ),
		// validation.Field(
		// 	&r.PrefectureID,
		// 	validation.Required.Error("募集エリアを選択してください"),
		// ),
		// validation.Field(
		// 	&r.CompetitionID,
		// 	validation.Required.Error("募集競技を選択してください"),
		// ),
	)
}

func (r *Recruitment) CreateRecruitment(ctx context.Context, client *ent.RecruitmentClient) (*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)

	res, err := client.
		Create().
		SetTitle(r.Title).
		SetType(recruitment.Type(strings.ToLower(string(r.Type)))).
		SetLevel(recruitment.Level(strings.ToLower(string(r.Level)))).
		SetNillableCapacity(r.Capacity).
		SetNillableStartAt(r.StartAt).
		SetNillableContent(r.Content).
		SetNillablePlace(r.Place).
		SetNillableLocationLat(r.LocationLat).
		SetNillableLocationLng(r.LocationLng).
		SetIsPublished(r.IsPublished).
		SetNillableClosingAt(r.ClosingAt).
		SetNillableCompetitionID(r.CompetitionID).
		SetNillablePrefectureID(r.PrefectureID).
		SetUserID(currentUser.ID).
		Save(ctx)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintln("recruitment create errors:", err.Error()))
		return &model.Recruitment{}, err
	}

	resRecruitment := &model.Recruitment{
		ID:          res.ID,
		Title:       res.Title,
		Type:        model.Type(res.Type),
		Level:       model.Level(res.Level),
		Place:       &res.Place,
		StartAt:     &res.StartAt,
		Content:     &res.Content,
		LocationLat: &res.LocationLat,
		LocationLng: &res.LocationLng,
		IsPublished: res.IsPublished,
		Capacity:    &res.Capacity,
		ClosingAt:   &res.ClosingAt,
		User:        currentUser,
	}

	return resRecruitment, nil
}
