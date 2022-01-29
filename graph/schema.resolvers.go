package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/domain/prefecture"
	"github.com/nagokos/connefut_backend/graph/domain/user"
	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (*model.User, error) {
	u := user.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	err := u.CreateUserValidate()

	if err != nil {
		logger.Log.Error().Msg(err.Error())
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(strings.ToLower(k), errMessage.Error()).AddGraphQLError(ctx)
		}

		return &model.User{}, err
	}

	res, err := u.CreateUser(r.client.User, ctx)
	if err != nil {
		return &model.User{}, err
	}

	token, _ := user.CreateToken(res.ID)

	auth.SetAuthCookie(ctx, token)

	return res, nil
}

func (r *mutationResolver) LoginUser(ctx context.Context, input model.LoginUserInput) (*model.User, error) {
	u := user.User{
		Email:    input.Email,
		Password: input.Password,
	}

	err := u.AuthenticateUserValidate()
	if err != nil {
		logger.Log.Error().Msg(err.Error())
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(strings.ToLower(k), errMessage.Error()).AddGraphQLError(ctx)
		}

		return &model.User{}, err
	}

	res, err := u.AuthenticateUser(r.client.User, ctx)
	if err != nil {
		return nil, err
	}

	token, _ := user.CreateToken(res.ID)

	auth.SetAuthCookie(ctx, token)

	return res, nil
}

func (r *mutationResolver) LogoutUser(ctx context.Context) (*bool, error) {
	auth.RemoveAuthCookie(ctx)
	return nil, nil
}

func (r *queryResolver) GetPrefectures(ctx context.Context) ([]*model.Prefecture, error) {
	res, err := prefecture.GetPrefectures(*r.client.Prefecture, ctx)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *queryResolver) GetCurrentUser(ctx context.Context) (*model.User, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, nil
	}
	return user, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
