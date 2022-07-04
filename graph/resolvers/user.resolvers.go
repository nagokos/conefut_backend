package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/user"
	"github.com/nagokos/connefut_backend/logger"
)

func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterUserInput) (*model.RegisterUserPayload, error) {
	u := user.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	err := u.CreateUserValidate()

	if err != nil {
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)

		var payload model.RegisterUserPayload

		for k, errMessage := range errs {
			payload.UserErrors = append(payload.UserErrors, &model.RegisterUserInvalidInputError{
				Message: errMessage.Error(),
				Field:   model.RegisterUserInvalidInputField(strings.ToLower(k)),
			})
		}

		return &payload, nil
	}

	payload, err := u.UserRegister(ctx, r.dbPool)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return payload, nil
	}

	token, _ := user.CreateToken(*payload.User.DatabaseID)
	auth.SetAuthCookie(ctx, token)

	return payload, nil
}

func (r *mutationResolver) LoginUser(ctx context.Context, input model.LoginUserInput) (*model.LoginUserPayload, error) {
	u := user.User{
		Email:    input.Email,
		Password: input.Password,
	}

	err := u.AuthenticateUserValidate()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)

		var payload model.LoginUserPayload

		for k, errMessage := range errs {
			payload.UserErrors = append(payload.UserErrors, model.LoginUserInvalidInputError{
				Message: errMessage.Error(),
				Field:   model.LoginUserInvalidInputField(strings.ToLower(k)),
			})
		}

		return &payload, nil
	}

	payload, err := u.UserLogin(ctx, r.dbPool)
	if err != nil {
		return payload, nil
	}

	token, _ := user.CreateToken(*payload.User.DatabaseID)
	auth.SetAuthCookie(ctx, token)

	return payload, nil
}

func (r *mutationResolver) LogoutUser(ctx context.Context) (bool, error) {
	auth.RemoveAuthCookie(ctx)
	return true, nil
}

func (r *queryResolver) CurrentUser(ctx context.Context) (*model.User, error) {
	user := auth.ForContext(ctx)
	return user, nil
}
