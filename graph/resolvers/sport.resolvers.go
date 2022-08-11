package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/sport"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

// Sports is the resolver for the sports field.
func (r *queryResolver) Sports(ctx context.Context) ([]*model.Sport, error) {
	sport, err := sport.GetSports(ctx, r.dbPool)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return sport, nil
}

// ID is the resolver for the id field.
func (r *sportResolver) ID(ctx context.Context, obj *model.Sport) (string, error) {
	return utils.GenerateUniqueID("Sport", obj.DatabaseID), nil
}

// Sport returns generated.SportResolver implementation.
func (r *Resolver) Sport() generated.SportResolver { return &sportResolver{r} }

type sportResolver struct{ *Resolver }
