package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/nagokos/connefut_backend/graph/domain/user"
	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/vektah/gqlparser/v2/gqlerror"
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

		errors := make(map[string]interface{})

		for k, errMessage := range errs {
			errors[strings.ToLower(k)] = fmt.Sprint(errMessage)
		}

		graphql.AddError(ctx, &gqlerror.Error{
			Extensions: errors,
		})
		return &model.User{}, err
	}

	res, err := u.CreateUser(r.client.User, ctx)

	if err != nil {
		graphql.AddError(ctx, &gqlerror.Error{
			Extensions: map[string]interface{}{
				"email": "このメールアドレスは既に存在します",
			},
		})
	}

	return res, err
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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
