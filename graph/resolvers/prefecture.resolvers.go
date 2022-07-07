package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/prefecture"
	"github.com/nagokos/connefut_backend/graph/utils"
)

// ID is the resolver for the id field.
func (r *prefectureResolver) ID(ctx context.Context, obj *model.Prefecture) (string, error) {
	return utils.GenerateUniqueID("Prefecture", obj.DatabaseID), nil
}

// Prefectures is the resolver for the prefectures field.
func (r *queryResolver) Prefectures(ctx context.Context) ([]*model.Prefecture, error) {
	res, err := prefecture.GetPrefectures(ctx, r.dbPool)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Prefecture returns generated.PrefectureResolver implementation.
func (r *Resolver) Prefecture() generated.PrefectureResolver { return &prefectureResolver{r} }

type prefectureResolver struct{ *Resolver }
