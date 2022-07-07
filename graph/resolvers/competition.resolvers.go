package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/competition"
	"github.com/nagokos/connefut_backend/graph/utils"
)

// ID is the resolver for the id field.
func (r *competitionResolver) ID(ctx context.Context, obj *model.Competition) (string, error) {
	return utils.GenerateUniqueID("Competition", obj.DatabaseID), nil
}

// Competitions is the resolver for the competitions field.
func (r *queryResolver) Competitions(ctx context.Context) ([]*model.Competition, error) {
	res, err := competition.GetCompetitions(ctx, r.dbPool)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Competition returns generated.CompetitionResolver implementation.
func (r *Resolver) Competition() generated.CompetitionResolver { return &competitionResolver{r} }

type competitionResolver struct{ *Resolver }
