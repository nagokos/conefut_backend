package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/domain/competition"
	"github.com/nagokos/connefut_backend/graph/domain/prefecture"
	"github.com/nagokos/connefut_backend/graph/domain/recruitment"
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

	if err != nil {
		logger.Log.Error().Msg(err.Error())
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(strings.ToLower(k), errMessage.Error()).AddGraphQLError(ctx)
		}

		return &model.User{}, err
	}

	res, err := u.CreateUser(r.client.User, ctx)
	if err != nil {
		return &model.User{}, err
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
			utils.NewValidationError(strings.ToLower(k), errMessage.Error()).AddGraphQLError(ctx)
		}

		return &model.User{}, err
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
		Level:         input.Level,
		Capacity:      input.Capacity,
		Place:         input.Place,
		LocationLat:   input.LocationLat,
		LocationLng:   input.LocationLng,
		IsPublished:   input.IsPublished,
		ClosingAt:     input.ClosingAt,
		CompetitionID: input.CompetitionID,
		PrefectureID:  input.PrefectureID,
	}

	err := rm.CreateRecruitmentValidate()
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintln("recruitment validation errors:", err.Error()))
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(strings.ToLower(k), errMessage.Error()).AddGraphQLError(ctx)
		}
		return &model.Recruitment{}, err
	}

	resRecruitment, err := rm.CreateRecruitment(ctx, r.client.Recruitment)
	if err != nil {
		return &model.Recruitment{}, err
	}

	return resRecruitment, nil
}

func (r *mutationResolver) DeleteRecruitment(ctx context.Context, id string) (bool, error) {
	res, err := recruitment.DeleteRecruitment(ctx, *r.client, id)
	if err != nil {
		return res, err
	}
	return res, nil
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

func (r *queryResolver) GetCurrentUserRecruitments(ctx context.Context) ([]*model.Recruitment, error) {
	res, err := recruitment.GetCurrentUserRecruitments(ctx, *r.client)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *queryResolver) GetEditRecruitment(ctx context.Context, id string) (*model.Recruitment, error) {
	res, err := recruitment.GetEditRecruitment(ctx, *r.client, id)
	if err != nil {
		return res, err
	}
	return res, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
