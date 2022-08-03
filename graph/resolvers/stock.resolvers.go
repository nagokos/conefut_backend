package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/stock"
	"github.com/nagokos/connefut_backend/graph/utils"
)

// AddStock is the resolver for the AddStock field.
func (r *mutationResolver) AddStock(ctx context.Context, recruitmentID string) (*model.FeedbackStock, error) {
	feedback, err := stock.AddStock(ctx, r.dbPool, utils.DecodeUniqueIDIdentifierOnly(recruitmentID))
	if err != nil {
		return nil, err
	}

	return feedback, nil
}

// RemoveStock is the resolver for the RemoveStock field.
func (r *mutationResolver) RemoveStock(ctx context.Context, recruitmentID string) (*model.FeedbackStock, error) {
	feedback, err := stock.RemoveStock(ctx, r.dbPool, utils.DecodeUniqueIDIdentifierOnly(recruitmentID))
	if err != nil {
		return nil, err
	}

	return feedback, nil
}

// CheckStocked is the resolver for the checkStocked field.
func (r *queryResolver) CheckStocked(ctx context.Context, recruitmentID string) (*model.FeedbackStock, error) {
	feedback, err := stock.CheckStocked(ctx, r.dbPool, utils.DecodeUniqueIDIdentifierOnly(recruitmentID))
	if err != nil {
		return nil, err
	}

	return feedback, nil
}

// GetStockedCount is the resolver for the getStockedCount field.
func (r *queryResolver) GetStockedCount(ctx context.Context, recruitmentID string) (*model.FeedbackStock, error) {
	panic(fmt.Errorf("not implemented"))
}
