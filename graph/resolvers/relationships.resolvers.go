package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/relationship"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

// ID is the resolver for the id field.
func (r *feedbackFollowResolver) ID(ctx context.Context, obj *model.FeedbackFollow) (string, error) {
	return utils.GenerateUniqueID("Relationship", obj.UserID), nil
}

// FollowingsCount is the resolver for the followingsCount field.
func (r *feedbackFollowResolver) FollowingsCount(ctx context.Context, obj *model.FeedbackFollow) (int, error) {
	count, err := relationship.GetFollowingsCount(ctx, r.dbPool, obj.UserID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return 0, err
	}
	return count, nil
}

// Follow is the resolver for the follow field.
func (r *mutationResolver) Follow(ctx context.Context, userID string) (*model.FollowResult, error) {
	feedback, err := relationship.Follow(ctx, r.dbPool, utils.DecodeUniqueIDIdentifierOnly(userID))
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return feedback, nil
}

// UnFollow is the resolver for the unFollow field.
func (r *mutationResolver) UnFollow(ctx context.Context, userID string) (*model.UnFollowResult, error) {
	feedback, err := relationship.UnFollow(ctx, r.dbPool, utils.DecodeUniqueIDIdentifierOnly(userID))
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return feedback, nil
}

// FeedbackFollow returns generated.FeedbackFollowResolver implementation.
func (r *Resolver) FeedbackFollow() generated.FeedbackFollowResolver {
	return &feedbackFollowResolver{r}
}

type feedbackFollowResolver struct{ *Resolver }
