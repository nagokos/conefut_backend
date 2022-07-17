package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/loader"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/recruitment"
	"github.com/nagokos/connefut_backend/graph/models/search"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

// CreateRecruitment is the resolver for the createRecruitment field.
func (r *mutationResolver) CreateRecruitment(ctx context.Context, input model.RecruitmentInput) (*model.RecruitmentPayload, error) {
	rm := recruitment.Recruitment{
		Title:         input.Title,
		Type:          input.Type,
		Detail:        input.Detail,
		StartAt:       input.StartAt,
		Venue:         input.Venue,
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
		return nil, err
	}

	payload, err := rm.CreateRecruitment(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// UpdateRecruitment is the resolver for the updateRecruitment field.
func (r *mutationResolver) UpdateRecruitment(ctx context.Context, id string, input model.RecruitmentInput) (*model.Recruitment, error) {
	viewer := auth.ForContext(ctx)
	if model.Status(input.Status) == model.StatusPublished &&
		viewer.EmailVerificationStatus == model.EmailVerificationStatusPending {
		return &model.Recruitment{}, errors.New("メールアドレスを認証してください")
	}

	rm := recruitment.Recruitment{
		Title:         input.Title,
		Type:          input.Type,
		Detail:        input.Detail,
		StartAt:       input.StartAt,
		Venue:         input.Venue,
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

// DeleteRecruitment is the resolver for the deleteRecruitment field.
func (r *mutationResolver) DeleteRecruitment(ctx context.Context, id string) (*model.DeleteRecruitmentPayload, error) {
	res, err := recruitment.DeleteRecruitment(ctx, r.dbPool, id)
	if err != nil {
		return res, err
	}
	return res, nil
}

// Recruitments is the resolver for the recruitments field.
func (r *queryResolver) Recruitments(ctx context.Context, first *int, after *string) (*model.RecruitmentConnection, error) {
	sp, err := search.NewSearchParams(first, after, nil, nil)
	if err != nil {
		return nil, err
	}

	res, err := recruitment.GetRecruitments(ctx, r.dbPool, sp)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ViewerRecruitments is the resolver for the viewerRecruitments field.
func (r *queryResolver) ViewerRecruitments(ctx context.Context, first *int, after *string) (*model.RecruitmentConnection, error) {
	params, err := search.NewSearchParams(first, after, nil, nil)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	res, err := recruitment.GetViewerRecruitments(ctx, r.dbPool, params)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Recruitment is the resolver for the recruitment field.
func (r *queryResolver) Recruitment(ctx context.Context, id string) (*model.Recruitment, error) {
	res, err := recruitment.GetRecruitment(ctx, r.dbPool, utils.DecodeUniqueID(id))
	if err != nil {
		return res, err
	}
	return res, nil
}

// StockedRecruitments is the resolver for the stockedRecruitments field.
func (r *queryResolver) StockedRecruitments(ctx context.Context) ([]*model.Recruitment, error) {
	res, err := recruitment.GetStockedRecruitments(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}

	return res, err
}

// AppliedRecruitments is the resolver for the appliedRecruitments field.
func (r *queryResolver) AppliedRecruitments(ctx context.Context) ([]*model.Recruitment, error) {
	res, err := recruitment.GetAppliedRecruitments(ctx, r.dbPool)
	if err != nil {
		return res, err
	}

	return res, err
}

// ID is the resolver for the id field.
func (r *recruitmentResolver) ID(ctx context.Context, obj *model.Recruitment) (string, error) {
	return utils.GenerateUniqueID("Recruitment", obj.DatabaseID), nil
}

// Competition is the resolver for the competition field.
func (r *recruitmentResolver) Competition(ctx context.Context, obj *model.Recruitment) (*model.Competition, error) {
	competition, err := loader.GetCompetition(ctx, obj.CompetitionID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return competition, nil
}

// Prefecture is the resolver for the prefecture field.
func (r *recruitmentResolver) Prefecture(ctx context.Context, obj *model.Recruitment) (*model.Prefecture, error) {
	if obj.PrefectureID == nil {
		return &model.Prefecture{}, nil
	}
	prefecture, err := loader.GetPrefecture(ctx, *obj.PrefectureID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return prefecture, nil
}

// User is the resolver for the user field.
func (r *recruitmentResolver) User(ctx context.Context, obj *model.Recruitment) (*model.User, error) {
	user, err := loader.GetUser(ctx, obj.UserID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return user, nil
}

// Tags is the resolver for the tags field.
func (r *recruitmentResolver) Tags(ctx context.Context, obj *model.Recruitment) ([]*model.Tag, error) {
	tags, err := loader.LoadTagsByRecruitmentID(ctx, obj.DatabaseID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return tags, nil
}

// Applicant is the resolver for the applicant field.
func (r *recruitmentResolver) Applicant(ctx context.Context, obj *model.Recruitment) (*model.Applicant, error) {
	panic(fmt.Errorf("not implemented"))
}

// Recruitment returns generated.RecruitmentResolver implementation.
func (r *Resolver) Recruitment() generated.RecruitmentResolver { return &recruitmentResolver{r} }

type recruitmentResolver struct{ *Resolver }
