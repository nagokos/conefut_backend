package recruitment

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/ent/recruitment"
	"github.com/nagokos/connefut_backend/ent/recruitmenttag"
	"github.com/nagokos/connefut_backend/ent/stock"
	"github.com/nagokos/connefut_backend/ent/user"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

type Recruitment struct {
	Title         string
	Type          model.Type
	Place         *string
	StartAt       *time.Time
	Content       *string
	LocationLat   *float64
	LocationLng   *float64
	Status        model.Status
	Capacity      *int
	ClosingAt     *time.Time
	CompetitionID *string
	PrefectureID  *string
	Tags          []*model.RecruitmentTagInput
}

func requiredIfUnnecessaryType() validation.RuleFunc {
	return func(v interface{}) error {
		if v == model.TypeUnnecessary {
			return errors.New("募集タイプを選択してください")
		}
		return nil
	}
}

func checkWithinTheDeadline(start time.Time) validation.RuleFunc {
	return func(v interface{}) error {
		var err error
		switch s := v.(type) {
		case *time.Time:
			difference := start.Sub(*s)
			if difference < 0 {
				err = errors.New("募集期限は開催日時よりも前に設定してください")
			} else {
				err = nil
			}
		}
		return err
	}
}

func beforeNowStart(v interface{}) error {
	var err error
	switch t := v.(type) {
	case *time.Time:
		difference := time.Since(*t).Minutes()
		if difference >= 1 {
			err = errors.New("開催日時は現在以降に設定してください")
		} else {
			err = nil
		}
	}
	return err
}

func beforeNowClosing(v interface{}) error {
	var err error
	switch t := v.(type) {
	case *time.Time:
		difference := time.Since(*t).Minutes()
		if difference >= 1 {
			err = errors.New("募集期限は現在以降に設定してください")
		} else {
			err = nil
		}
	}
	return err
}

func (r Recruitment) RecruitmentValidate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.Title,
			validation.Required.Error("タイトルを入力してください"),
			validation.RuneLength(1, 60).Error("タイトルは1文字以上60文字以内で入力してください"),
		),
		validation.Field(
			&r.CompetitionID,
			validation.When(r.Status == model.StatusPublished,
				validation.Required.Error("募集競技を選択してください"),
			),
		),
		validation.Field(
			&r.Type,
			validation.In(
				model.TypeUnnecessary,
				model.TypeOpponent,
				model.TypeIndividual,
				model.TypeMember,
				model.TypeJoining,
				model.TypeOthers,
			),
			validation.When(r.Status == model.StatusPublished,
				validation.By(requiredIfUnnecessaryType()),
			),
		),
		validation.Field(
			&r.Content,
			validation.When(r.Status == model.StatusPublished,
				validation.Required.Error("募集の詳細を入力してください"),
				validation.RuneLength(1, 10000).Error("募集の詳細は10000文字以内で入力してください"),
			).Else(validation.RuneLength(0, 10000).Error("募集の詳細は10000文字以内で入力してください")),
		),
		validation.Field(
			&r.PrefectureID,
			validation.When(r.Status == model.StatusPublished,
				validation.Required.Error("募集エリアを選択してください"),
			),
		),
		validation.Field(
			&r.Place,
			validation.When(r.Status == model.StatusPublished,
				validation.When(r.Type == model.TypeOpponent || r.Type == model.TypeIndividual,
					validation.Required.Error("会場名を入力してください"),
				),
			),
		),
		validation.Field(
			&r.Capacity,
			validation.When(r.Status == model.StatusPublished,
				validation.When(
					r.Type == model.TypeOpponent ||
						r.Type == model.TypeIndividual,
					validation.Required.Error("募集人数は1名以上にしてください"),
					validation.Min(1).Error("募集人数は1名以上にしてください"),
				),
			),
		),
		validation.Field(
			&r.StartAt,
			validation.When(r.Status == model.StatusPublished,
				validation.When(r.Type == model.TypeOpponent || r.Type == model.TypeIndividual,
					validation.By(beforeNowStart),
					validation.Required.Error("開催日時を設定してください"),
				),
			),
		),
		validation.Field(
			&r.ClosingAt,
			validation.When(r.Status == model.StatusPublished,
				validation.Required.Error("募集期限を設定してください"),
				validation.When(r.Type == model.TypeOpponent || r.Type == model.TypeIndividual,
					validation.By(beforeNowClosing),
					validation.By(checkWithinTheDeadline(*r.StartAt)),
				),
			),
		),
	)
}

func (r *Recruitment) CreateRecruitment(ctx context.Context, client *ent.RecruitmentClient) (*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return &model.Recruitment{}, errors.New("ログインしてください")
	}

	res, err := client.
		Create().
		SetTitle(r.Title).
		SetType(recruitment.Type(strings.ToLower(string(r.Type)))).
		SetNillableCapacity(r.Capacity).
		SetNillableStartAt(r.StartAt).
		SetNillableContent(r.Content).
		SetNillablePlace(r.Place).
		SetNillableLocationLat(r.LocationLat).
		SetNillableLocationLng(r.LocationLng).
		SetStatus(recruitment.Status(strings.ToLower(string(r.Status)))).
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
		Place:       &res.Place,
		StartAt:     &res.StartAt,
		Content:     &res.Content,
		LocationLat: &res.LocationLat,
		LocationLng: &res.LocationLng,
		Status:      model.Status(res.Status),
		Capacity:    &res.Capacity,
		ClosingAt:   &res.ClosingAt,
		User:        &model.User{},
	}

	fmt.Println(*resRecruitment)

	return resRecruitment, nil
}

func GetCurrentUserRecruitments(ctx context.Context, client ent.Client) ([]*model.Recruitment, error) {
	var recruitments []*model.Recruitment

	currentUser := auth.ForContext(ctx)

	res, err := client.User.
		Query().
		Where(
			user.ID(currentUser.ID),
		).
		QueryRecruitments().
		WithCompetition().
		Order(
			ent.Desc(recruitment.FieldStatus),
			ent.Desc(recruitment.FieldUpdatedAt),
		).
		All(ctx)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("get currentUser recruitment error %s", err.Error()))
		return recruitments, err
	}

	for _, recruitment := range res {
		var compName string
		if recruitment.Edges.Competition != nil {
			compName = recruitment.Edges.Competition.Name
		}
		recruitments = append(recruitments, &model.Recruitment{
			ID:     recruitment.ID,
			Title:  recruitment.Title,
			Type:   model.Type(strings.ToUpper(string(recruitment.Type))),
			Status: model.Status(strings.ToUpper(string(recruitment.Status))),
			Competition: &model.Competition{
				Name: compName,
			},
		})
	}

	return recruitments, nil
}

func GetRecruitment(ctx context.Context, client ent.Client, id string) (*model.Recruitment, error) {
	var prefecture *model.Prefecture
	var competition *model.Competition
	var user *model.User

	res, err := client.Recruitment.
		Query().
		Where(recruitment.ID(id)).
		WithCompetition().
		WithPrefecture().
		WithUser().
		Only(ctx)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("get recruitment error %s", err.Error()))
		return nil, err
	}

	if res.Edges.Prefecture != nil {
		prefecture = &model.Prefecture{
			ID:   res.Edges.Prefecture.ID,
			Name: res.Edges.Prefecture.Name,
		}
	}

	if res.Edges.Competition != nil {
		competition = &model.Competition{
			ID:   res.Edges.Competition.ID,
			Name: res.Edges.Competition.Name,
		}
	}

	if res.Edges.User != nil {
		user = &model.User{
			ID:     res.Edges.User.ID,
			Name:   res.Edges.User.Name,
			Avatar: res.Edges.User.Avatar,
		}
	}

	resRecruitment := &model.Recruitment{
		ID:          res.ID,
		Title:       res.Title,
		Type:        model.Type(strings.ToUpper(string(res.Type))),
		Place:       &res.Place,
		StartAt:     &res.StartAt,
		Content:     &res.Content,
		LocationLat: &res.LocationLat,
		LocationLng: &res.LocationLng,
		Status:      model.Status(strings.ToUpper(string(res.Status))),
		Capacity:    &res.Capacity,
		ClosingAt:   &res.ClosingAt,
		Prefecture:  prefecture,
		Competition: competition,
		UpdatedAt:   res.UpdatedAt,
		User:        user,
	}
	return resRecruitment, nil
}

func GetRecruitments(ctx context.Context, client ent.Client) ([]*model.Recruitment, error) {
	var resRecruitments []*model.Recruitment
	res, err := client.Recruitment.
		Query().
		Where(
			recruitment.StatusEQ(recruitment.StatusPublished),
			recruitment.ClosingAtGT(time.Now().Local()),
		).
		WithCompetition().
		WithPrefecture().
		WithUser().
		All(ctx)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("get recruitments error %s", err.Error()))
		return []*model.Recruitment{}, err
	}

	for _, recruitment := range res {
		resRecruitments = append(resRecruitments, &model.Recruitment{
			ID:        recruitment.ID,
			Title:     recruitment.Title,
			Type:      model.Type(strings.ToUpper(string(recruitment.Type))),
			Content:   &recruitment.Content,
			StartAt:   &recruitment.StartAt,
			UpdatedAt: recruitment.UpdatedAt,
			ClosingAt: &recruitment.ClosingAt,
			Capacity:  &recruitment.Capacity,
			Place:     &recruitment.Place,
			User: &model.User{
				Name:   recruitment.Edges.User.Name,
				Avatar: recruitment.Edges.User.Avatar,
			},
			Prefecture: &model.Prefecture{
				Name: recruitment.Edges.Prefecture.Name,
			},
			Status: model.Status(strings.ToUpper(string(recruitment.Status))),
		})
	}

	fmt.Println(resRecruitments)
	return resRecruitments, nil
}

func GetStockedRecruitments(ctx context.Context, client ent.Client) ([]*model.Recruitment, error) {
	var recruitments []*model.Recruitment

	currentUser := auth.ForContext(ctx)

	res, err := client.Recruitment.
		Query().
		WithUser().
		Where(
			recruitment.HasStocksWith(
				stock.UserID(currentUser.ID),
			),
		).
		Order(ent.Desc(recruitment.FieldStatus)).
		All(ctx)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("get stocked recruitments error %s", err.Error()))
		return []*model.Recruitment{}, err
	}

	for _, recruitment := range res {
		var user *ent.User
		if recruitment.Edges.User != nil {
			user = recruitment.Edges.User
		}

		recruitments = append(recruitments, &model.Recruitment{
			ID:     recruitment.ID,
			Title:  recruitment.Title,
			Type:   model.Type(strings.ToUpper(string(recruitment.Type))),
			Status: model.Status(strings.ToUpper(string(recruitment.Status))),
			User: &model.User{
				ID:     user.ID,
				Name:   user.Name,
				Avatar: user.Avatar,
			},
		})
	}

	return recruitments, nil
}

func (r *Recruitment) UpdateRecruitment(ctx context.Context, client ent.Client, id string) (*model.Recruitment, error) {
	var entRecruitment *ent.Recruitment
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return nil, errors.New("ログインしてください")
	}

	i, err := client.Recruitment.
		Update().
		Where(
			recruitment.ID(id),
			recruitment.HasUserWith(user.ID(currentUser.ID)),
		).
		ClearCapacity().
		ClearLocationLat().
		ClearLocationLng().
		SetTitle(r.Title).
		SetType(recruitment.Type(strings.ToLower(string(r.Type)))).
		SetNillableCapacity(r.Capacity).
		SetNillableStartAt(r.StartAt).
		SetNillableContent(r.Content).
		SetNillablePlace(r.Place).
		SetNillableLocationLat(r.LocationLat).
		SetNillableLocationLng(r.LocationLng).
		SetStatus(recruitment.Status(strings.ToLower(string(r.Status)))).
		SetNillableClosingAt(r.ClosingAt).
		SetNillableCompetitionID(r.CompetitionID).
		SetNillablePrefectureID(r.PrefectureID).
		Save(ctx)
	if i == 0 {
		return nil, errors.New("募集の更新に失敗しました")
	}
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("recruitment update error %s", err.Error()))
		return nil, err
	}

	entRecruitment, err = client.Recruitment.Query().Where(recruitment.ID(id)).Only(ctx)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("recruitment update error %s", err.Error()))
		return nil, err
	}

	resRecruitment := &model.Recruitment{
		ID:     entRecruitment.ID,
		Title:  entRecruitment.Title,
		Status: model.Status(strings.ToUpper(string(entRecruitment.Status))),
	}
	return resRecruitment, nil
}

func DeleteRecruitment(ctx context.Context, client ent.Client, id string) (bool, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return false, errors.New("ログインしてください")
	}

	res, err := client.Recruitment.
		Delete().
		Where(
			recruitment.HasUserWith(user.ID(currentUser.ID)),
			recruitment.ID(id),
		).
		Exec(ctx)
	if res == 0 {
		return false, errors.New("募集の削除に失敗しました")
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
