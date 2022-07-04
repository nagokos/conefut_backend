package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/tag"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

func (r *mutationResolver) CreateTag(ctx context.Context, input model.CreateTagInput) (*model.Tag, error) {
	tag := tag.Tag{
		Name: input.Name,
	}

	err := tag.CreateTagValidate()
	if err != nil {
		logger.NewLogger().Sugar().Errorf("recruitment validation errors %s", err.Error())
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(errMessage.Error(), utils.WithField(strings.ToLower(k))).AddGraphQLError(ctx)
		}
		return nil, errors.New("タグの作成に失敗しました")
	}

	res, err := tag.CreateTag(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *queryResolver) Tags(ctx context.Context) ([]*model.Tag, error) {
	res, err := tag.GetTags(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}
	return res, nil
}