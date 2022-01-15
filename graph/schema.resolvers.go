package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/nagokos/connefut_backend/graph/domain/prefecture"
	"github.com/nagokos/connefut_backend/graph/domain/user"
	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	u := user.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	err := u.Validate()

	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintln(err))
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			NewValidationError(strings.ToLower(k), errMessage.Error()).AddGraphQLError(ctx)
		}

		return &model.User{}, nil
	}

	res, err := u.CreateUser(r.client.User, ctx)

	if err != nil {
		NewValidationError("email", "このメールアドレスは既に使用されています").AddGraphQLError(ctx)
	}

	return res, nil
}

func (r *mutationResolver) UpdateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) User(ctx context.Context, userID *string) (*model.User, error) {
	var resUser model.User

	res, err := r.client.User.Get(ctx, *userID)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintln(err))
		return &resUser, fmt.Errorf("user not found")
	}
	resUser = model.User{
		ID:    res.ID,
		Name:  res.Name,
		Email: res.Email,
	}
	return &resUser, nil
}

func (r *queryResolver) GetPrefectures(ctx context.Context) ([]*model.Prefecture, error) {
	res, err := prefecture.GetPrefectures(*r.client.Prefecture, ctx)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
