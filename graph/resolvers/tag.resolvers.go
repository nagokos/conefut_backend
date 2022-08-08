package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/tag"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

// CreateTag is the resolver for the createTag field.
func (r *mutationResolver) CreateTag(ctx context.Context, input model.CreateTagInput) (model.CreateTagResult, error) {
	tag := tag.Tag{
		Name: input.Name,
	}

	if err := tag.CreateTagValidate(); err != nil {
		logger.NewLogger().Sugar().Errorf("recruitment validation errors %s", err.Error())
		errs := err.(validation.Errors)

		var result model.CreateTagInvalidInputErrors
		for field, message := range errs {
			result.InvalidInputs = append(result.InvalidInputs, &model.CreateTagInvalidInputError{
				Field:   model.CreateTagInvalidInputField(strings.ToUpper(field)),
				Message: message.Error(),
			})
		}
		return result, nil
	}

	result, err := tag.CreateTag(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Tags is the resolver for the tags field.
func (r *queryResolver) Tags(ctx context.Context, first int) (*model.TagConnection, error) {
	payload, err := tag.GetTags(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}
	return payload, nil
}

// ID is the resolver for the id field.
func (r *tagResolver) ID(ctx context.Context, obj *model.Tag) (string, error) {
	return utils.GenerateUniqueID("Tag", obj.DatabaseID), nil
}

// Cursor is the resolver for the cursor field.
func (r *tagEdgeResolver) Cursor(ctx context.Context, obj *model.TagEdge) (string, error) {
	return utils.GenerateUniqueID("Tag", obj.Node.DatabaseID), nil
}

// Tag returns generated.TagResolver implementation.
func (r *Resolver) Tag() generated.TagResolver { return &tagResolver{r} }

// TagEdge returns generated.TagEdgeResolver implementation.
func (r *Resolver) TagEdge() generated.TagEdgeResolver { return &tagEdgeResolver{r} }

type tagResolver struct{ *Resolver }
type tagEdgeResolver struct{ *Resolver }
