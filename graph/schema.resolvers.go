package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/domain/competition"
	"github.com/nagokos/connefut_backend/graph/domain/prefecture"
	"github.com/nagokos/connefut_backend/graph/domain/recruitment"
	"github.com/nagokos/connefut_backend/graph/domain/stock"
	"github.com/nagokos/connefut_backend/graph/domain/tag"
	"github.com/nagokos/connefut_backend/graph/domain/user"
	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	u := user.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	err := u.CreateUserValidate()

	fmt.Println(err)

	if err != nil {
		logger.Log.Error().Msg(err.Error())
		errs := err.(validation.Errors)

		fmt.Println(errs)

		for k, errMessage := range errs {
			utils.NewValidationError(errMessage.Error(), utils.WithField(strings.ToLower(k))).AddGraphQLError(ctx)
		}

		return nil, errors.New("フォームに不備があります")
	}

	res, err := u.CreateUser(r.client.User, ctx)
	if err != nil {
		return nil, err
	}

	token, _ := user.CreateToken(res.ID)

	auth.SetAuthCookie(ctx, token)

	return res, nil
}

func (r *mutationResolver) LoginUser(ctx context.Context, input model.LoginUserInput) (*model.User, error) {
	u := user.User{
		Email:    input.Email,
		Password: input.Password,
	}

	err := u.AuthenticateUserValidate()
	if err != nil {
		logger.Log.Error().Msg(err.Error())
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(errMessage.Error(), utils.WithField(strings.ToLower(k))).AddGraphQLError(ctx)
		}

		return nil, errors.New("フォームに不備があります")
	}

	res, err := u.AuthenticateUser(r.client.User, ctx)
	if err != nil {
		return nil, err
	}

	token, _ := user.CreateToken(res.ID)

	auth.SetAuthCookie(ctx, token)

	return res, nil
}

func (r *mutationResolver) LogoutUser(ctx context.Context) (bool, error) {
	auth.RemoveAuthCookie(ctx)
	return true, nil
}

func (r *mutationResolver) CreateRecruitment(ctx context.Context, input model.RecruitmentInput) (*model.Recruitment, error) {
	rm := recruitment.Recruitment{
		Title:         input.Title,
		Type:          input.Type,
		Content:       input.Content,
		StartAt:       input.StartAt,
		Capacity:      input.Capacity,
		Place:         input.Place,
		LocationLat:   input.LocationLat,
		LocationLng:   input.LocationLng,
		Status:        input.Status,
		ClosingAt:     input.ClosingAt,
		CompetitionID: input.CompetitionID,
		PrefectureID:  input.PrefectureID,
		Tags:          input.Tags,
	}

	err := rm.RecruitmentValidate()
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintln("recruitment validation errors:", err.Error()))
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(errMessage.Error(), utils.WithField(strings.ToLower(k))).AddGraphQLError(ctx)
		}
		return &model.Recruitment{}, err
	}

	resRecruitment, err := rm.CreateRecruitment(ctx, r.client)
	if err != nil {
		return &model.Recruitment{}, err
	}

	return resRecruitment, nil
}

func (r *mutationResolver) UpdateRecruitment(ctx context.Context, id string, input model.RecruitmentInput) (*model.Recruitment, error) {
	rm := recruitment.Recruitment{
		Title:         input.Title,
		Type:          input.Type,
		Content:       input.Content,
		StartAt:       input.StartAt,
		Capacity:      input.Capacity,
		Place:         input.Place,
		LocationLat:   input.LocationLat,
		LocationLng:   input.LocationLng,
		Status:        input.Status,
		ClosingAt:     input.ClosingAt,
		CompetitionID: input.CompetitionID,
		PrefectureID:  input.PrefectureID,
		Tags:          input.Tags,
	}

	err := rm.RecruitmentValidate()
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("recruitment validation errors %s", err.Error()))
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(errMessage.Error(), utils.WithField(strings.ToLower(k))).AddGraphQLError(ctx)
		}
		return &model.Recruitment{}, err
	}

	res, err := rm.UpdateRecruitment(ctx, *r.client, id)
	if err != nil {
		logger.Log.Error().Msg(err.Error())
		return res, err
	}

	return res, nil
}

func (r *mutationResolver) DeleteRecruitment(ctx context.Context, id string) (bool, error) {
	res, err := recruitment.DeleteRecruitment(ctx, *r.client, id)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (r *mutationResolver) CreateStock(ctx context.Context, recruitmentID string) (bool, error) {
	_, err := stock.CreateStock(ctx, *r.client, recruitmentID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) DeleteStock(ctx context.Context, recruitmentID string) (bool, error) {
	_, err := stock.DeleteStock(ctx, *r.client, recruitmentID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) CreateTag(ctx context.Context, input model.CreateTagInput) (*model.Tag, error) {
	tag := tag.Tag{
		Name: input.Name,
	}

	err := tag.CreateTagValidate()
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("recruitment validation errors %s", err.Error()))
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(errMessage.Error(), utils.WithField(strings.ToLower(k))).AddGraphQLError(ctx)
		}
		return nil, errors.New("タグの作成に失敗しました")
	}

	res, err := tag.CreateTag(ctx, r.client)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *mutationResolver) AddRecruitmentTag(ctx context.Context, tagID string, recruitmentID string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetPrefectures(ctx context.Context) ([]*model.Prefecture, error) {
	res, err := prefecture.GetPrefectures(*r.client.Prefecture, ctx)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *queryResolver) GetCurrentUser(ctx context.Context) (*model.User, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, nil
	}
	return user, nil
}

func (r *queryResolver) GetCompetitions(ctx context.Context) ([]*model.Competition, error) {
	res, err := competition.GetCompetitions(ctx, r.client.Competition)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *queryResolver) GetRecruitments(ctx context.Context) ([]*model.Recruitment, error) {
	res, err := recruitment.GetRecruitments(ctx, *r.client)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (r *queryResolver) GetCurrentUserRecruitments(ctx context.Context) ([]*model.Recruitment, error) {
	res, err := recruitment.GetCurrentUserRecruitments(ctx, *r.client)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *queryResolver) GetRecruitment(ctx context.Context, id string) (*model.Recruitment, error) {
	res, err := recruitment.GetRecruitment(ctx, *r.client, id)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (r *queryResolver) GetStockedRecruitments(ctx context.Context) ([]*model.Recruitment, error) {
	res, err := recruitment.GetStockedRecruitments(ctx, *r.client)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (r *queryResolver) CheckStocked(ctx context.Context, recruitmentID string) (bool, error) {
	res, err := stock.CheckStocked(ctx, *r.client, recruitmentID)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (r *queryResolver) GetStockedCount(ctx context.Context, recruitmentID string) (int, error) {
	res, err := stock.GetStockedCount(ctx, *r.client, recruitmentID)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *queryResolver) GetTags(ctx context.Context) ([]*model.Tag, error) {
	res, err := tag.GetTags(ctx, r.client)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *queryResolver) GetRecruitmentTags(ctx context.Context, recruitmentID string) ([]*model.Tag, error) {
	res, err := tag.GetRecruitmentTags(ctx, r.client, recruitmentID)
	if err != nil {
		return res, err
	}
	return res, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
