package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/recruitment"
	"github.com/nagokos/connefut_backend/graph/models/search"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

func (r *mutationResolver) CreateRecruitment(ctx context.Context, input model.RecruitmentInput) (*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)
	if model.Status(input.Status) == model.StatusPublished &&
		currentUser.EmailVerificationStatus == model.EmailVerificationStatusPending {
		return &model.Recruitment{}, errors.New("メールアドレスを認証してください")
	}

	rm := recruitment.Recruitment{
		Title:         input.Title,
		Type:          input.Type,
		Detail:        input.Detail,
		StartAt:       input.StartAt,
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
		logger.NewLogger().Sugar().Errorf("recruitment validation errors:", err.Error())
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(errMessage.Error(), utils.WithField(strings.ToLower(k))).AddGraphQLError(ctx)
		}
		return &model.Recruitment{}, err
	}

	resRecruitment, err := rm.CreateRecruitment(ctx, r.dbPool)
	if err != nil {
		return &model.Recruitment{}, err
	}

	return resRecruitment, nil
}

func (r *mutationResolver) UpdateRecruitment(ctx context.Context, id string, input model.RecruitmentInput) (*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)
	if model.Status(input.Status) == model.StatusPublished &&
		currentUser.EmailVerificationStatus == model.EmailVerificationStatusPending {
		return &model.Recruitment{}, errors.New("メールアドレスを認証してください")
	}

	rm := recruitment.Recruitment{
		Title:         input.Title,
		Type:          input.Type,
		Detail:        input.Detail,
		StartAt:       input.StartAt,
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
		logger.NewLogger().Sugar().Errorf("recruitment validation errors %s", err.Error())
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(errMessage.Error(), utils.WithField(strings.ToLower(k))).AddGraphQLError(ctx)
		}
		return &model.Recruitment{}, err
	}

	res, err := rm.UpdateRecruitment(ctx, r.dbPool, id)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return res, err
	}

	return res, nil
}

func (r *mutationResolver) DeleteRecruitment(ctx context.Context, id string) (*model.Recruitment, error) {
	res, err := recruitment.DeleteRecruitment(ctx, r.dbPool, id)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (r *queryResolver) Recruitments(ctx context.Context, first *int, after *string, last *int, before *string) (*model.RecruitmentConnection, error) {
	sp, err := search.NewSearchParams(first, after, last, before)
	if err != nil {
		return nil, err
	}

	res, err := recruitment.GetRecruitments(ctx, r.dbPool, sp)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *queryResolver) CurrentUserRecruitments(ctx context.Context) ([]*model.Recruitment, error) {
	res, err := recruitment.GetCurrentUserRecruitments(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *queryResolver) Recruitment(ctx context.Context, id string) (*model.Recruitment, error) {
	res, err := recruitment.GetRecruitment(ctx, r.dbPool, id)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (r *queryResolver) StockedRecruitments(ctx context.Context) ([]*model.Recruitment, error) {
	res, err := recruitment.GetStockedRecruitments(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (r *queryResolver) AppliedRecruitments(ctx context.Context) ([]*model.Recruitment, error) {
	res, err := recruitment.GetAppliedRecruitments(ctx, r.dbPool)
	if err != nil {
		return res, err
	}

	return res, err
}
