package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/user"
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
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(errMessage.Error(), utils.WithField(strings.ToLower(k))).AddGraphQLError(ctx)
		}

		return nil, errors.New("フォームに不備があります")
	}

	resUser, err := u.CreateUser(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}

	token, _ := user.CreateToken(*resUser.DatabaseID)

	auth.SetAuthCookie(ctx, token)

	return resUser, nil
}

func (r *mutationResolver) LoginUser(ctx context.Context, input model.LoginUserInput) (*model.User, error) {
	u := user.User{
		Email:    input.Email,
		Password: input.Password,
	}

	err := u.AuthenticateUserValidate()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(errMessage.Error(), utils.WithField(strings.ToLower(k))).AddGraphQLError(ctx)
		}

		return nil, errors.New("フォームに不備があります")
	}

	resUser, err := u.Authenticate(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}

	token, _ := user.CreateToken(*resUser.DatabaseID)

	auth.SetAuthCookie(ctx, token)

	return nil, nil
}

func (r *mutationResolver) LogoutUser(ctx context.Context) (bool, error) {
	auth.RemoveAuthCookie(ctx)
	return true, nil
}

func (r *queryResolver) GetCurrentUser(ctx context.Context) (*model.User, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, nil
	}
	return user, nil
}
