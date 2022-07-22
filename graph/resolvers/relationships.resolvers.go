package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/nagokos/connefut_backend/graph/model"
)

// Follow is the resolver for the follow field.
func (r *mutationResolver) Follow(ctx context.Context, userID string) (*model.FeedbackFollow, error) {
	panic(fmt.Errorf("not implemented"))
}

// UnFollow is the resolver for the unFollow field.
func (r *mutationResolver) UnFollow(ctx context.Context, userID string) (*model.FeedbackFollow, error) {
	panic(fmt.Errorf("not implemented"))
}

// CheckFollowed is the resolver for the checkFollowed field.
func (r *queryResolver) CheckFollowed(ctx context.Context, userID string) (*model.FeedbackFollow, error) {
	panic(fmt.Errorf("not implemented"))
}
