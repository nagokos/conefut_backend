package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/competition"
)

func (r *queryResolver) GetCompetitions(ctx context.Context) ([]*model.Competition, error) {
	res, err := competition.GetCompetitions(ctx, r.dbPool)
	if err != nil {
		return res, err
	}

	return res, nil
}
