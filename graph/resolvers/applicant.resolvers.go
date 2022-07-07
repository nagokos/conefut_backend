package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/applicant"
	"github.com/nagokos/connefut_backend/graph/utils"
)

// ID is the resolver for the id field.
func (r *applicantResolver) ID(ctx context.Context, obj *model.Applicant) (string, error) {
	return utils.GenerateUniqueID("Applicant", obj.DatabaseID), nil
}

// ApplyForRecruitment is the resolver for the applyForRecruitment field.
func (r *mutationResolver) ApplyForRecruitment(ctx context.Context, recruitmentID string, input *model.ApplicantInput) (*model.ApplyForRecruitmentPayload, error) {
	_, err := applicant.CreateApplicant(ctx, r.dbPool, recruitmentID, input.Message)
	if err != nil {
		return nil, err
	}

	return nil, err
}

// CheckAppliedForRecruitment is the resolver for the checkAppliedForRecruitment field.
func (r *queryResolver) CheckAppliedForRecruitment(ctx context.Context, recruitmentID string) (*model.FeedbackApplicant, error) {
	_, err := applicant.CheckAppliedForRecruitment(ctx, r.dbPool, recruitmentID)
	if err != nil {
		return nil, err
	}
	return nil, err
}

// Applicant returns generated.ApplicantResolver implementation.
func (r *Resolver) Applicant() generated.ApplicantResolver { return &applicantResolver{r} }

type applicantResolver struct{ *Resolver }
