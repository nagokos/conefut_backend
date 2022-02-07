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
	Content       string
	LocationURL   *string
	Capacity      *int
	ClosingAt     time.Time
	CompetitionID string
	PrefectureID  string
}

func (r Recruitment) CreateRecruitmentValidate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.Title,
			validation.Required.Error("タイトルを入力してください"),
			validation.RuneLength(1, 60).Error("エラーです"),
		),
		validation.Field(
			&r.Type,
			validation.Required.Error("募集タイプを選択してください"),
		),
		validation.Field(
			&r.Content,
			validation.Required.Error("募集の詳細を入力してください"),
			validation.RuneLength(1, 10000).Error("募集の詳細は10000文字以内で入力してください"),
		),
		validation.Field(
			&r.PrefectureID,
			validation.Required.Error("募集エリアを選択してください"),
		),
		validation.Field(
			&r.CompetitionID,
			validation.Required.Error("募集競技を選択してください"),
		),
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
		SetContent(r.Content).
		SetNillablePlace(r.Place).
		SetNillableLocationURL(r.LocationURL).
		SetClosingAt(r.ClosingAt).
		SetCompetitionID(r.CompetitionID).
		SetPrefectureID(r.PrefectureID).
		SetUserID(currentUser.ID).
		Save(ctx)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintln("recruitment create errors:", err.Error()))
		return &model.Recruitment{}, err
	}

	c, _ := res.Competition(ctx)
	p, _ := res.Prefecture(ctx)

	resRecruitment := &model.Recruitment{
		ID:          res.ID,
		Title:       res.Title,
		Type:        model.Type(res.Type),
		Level:       model.Level(res.Level),
		Place:       &res.Place,
		StartAt:     &res.StartAt,
		Content:     res.Content,
		LocationURL: &res.LocationURL,
		Capacity:    &res.Capacity,
		ClosingAt:   res.ClosingAt,
		Competition: &model.Competition{
			ID:   c.ID,
			Name: c.Name,
		},
		Prefecture: &model.Prefecture{
			ID:   p.ID,
			Name: p.Name,
		},
		User: currentUser,
	}

	return resRecruitment, nil
}
