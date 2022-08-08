package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/stock"
	"github.com/nagokos/connefut_backend/graph/utils"
)

// ID is the resolver for the id field.
func (r *feedbackStockResolver) ID(ctx context.Context, obj *model.FeedbackStock) (string, error) {
	return utils.GenerateUniqueID("Stock", obj.RecruitmentID), nil
}

// AddStock is the resolver for the AddStock field.
func (r *mutationResolver) AddStock(ctx context.Context, recruitmentID string) (*model.AddStockResult, error) {
	feedback, err := stock.AddStock(ctx, r.dbPool, utils.DecodeUniqueIDIdentifierOnly(recruitmentID))
	if err != nil {
		return nil, err
	}

	return feedback, nil
}

// RemoveStock is the resolver for the RemoveStock field.
func (r *mutationResolver) RemoveStock(ctx context.Context, recruitmentID string) (*model.RemoveStockResult, error) {
	feedback, err := stock.RemoveStock(ctx, r.dbPool, utils.DecodeUniqueIDIdentifierOnly(recruitmentID))
	if err != nil {
		return nil, err
	}

	return feedback, nil
}

// RemovedRecruitmentID is the resolver for the removedRecruitmentId field.
func (r *removeStockResultResolver) RemovedRecruitmentID(ctx context.Context, obj *model.RemoveStockResult) (string, error) {
	return utils.GenerateUniqueID("Recruitment", obj.FeedbackStock.RecruitmentID), nil
}

// FeedbackStock returns generated.FeedbackStockResolver implementation.
func (r *Resolver) FeedbackStock() generated.FeedbackStockResolver { return &feedbackStockResolver{r} }

// RemoveStockResult returns generated.RemoveStockResultResolver implementation.
func (r *Resolver) RemoveStockResult() generated.RemoveStockResultResolver {
	return &removeStockResultResolver{r}
}

type feedbackStockResolver struct{ *Resolver }
type removeStockResultResolver struct{ *Resolver }
