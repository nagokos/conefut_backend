package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/relationship"
	"github.com/nagokos/connefut_backend/logger"
)

// Follow is the resolver for the follow field.
func (r *mutationResolver) Follow(ctx context.Context, userID string) (*model.FeedbackFollow, error) {
	feedback, err := relationship.Follow(ctx, r.dbPool, userID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return feedback, nil
}

// UnFollow is the resolver for the unFollow field.
func (r *mutationResolver) UnFollow(ctx context.Context, userID string) (*model.FeedbackFollow, error) {
	feedback, err := relationship.UnFollow(ctx, r.dbPool, userID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return feedback, nil
}

// CheckFollowed is the resolver for the checkFollowed field.
func (r *queryResolver) CheckFollowed(ctx context.Context, userID string) (*model.FeedbackFollow, error) {
	feedback, err := relationship.CheckFollowed(ctx, r.dbPool, userID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return feedback, nil
}

// CheckFollowedByRecruitmentID is the resolver for the checkFollowedByRecruitmentId field.
func (r *queryResolver) CheckFollowedByRecruitmentID(ctx context.Context, recruitmentID string) (*model.FeedbackFollow, error) {
	feedback, err := relationship.CheckFollowedByRecruitmentID(ctx, r.dbPool, recruitmentID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return feedback, nil
}
