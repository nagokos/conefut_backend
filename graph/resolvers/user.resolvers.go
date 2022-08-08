package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/graph/cookie"
	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/recruitment"
	"github.com/nagokos/connefut_backend/graph/models/relationship"
	"github.com/nagokos/connefut_backend/graph/models/search"
	"github.com/nagokos/connefut_backend/graph/models/user"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

// RegisterUser is the resolver for the registerUser field.
func (r *mutationResolver) RegisterUser(ctx context.Context, input model.RegisterUserInput) (model.RegisterUserResult, error) {
	u := user.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}

	if err := u.CreateUserValidate(); err != nil {
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)
		var result model.RegisterUserInvalidInputErrors
		for k, errMessage := range errs {
			result.InvalidInputs = append(result.InvalidInputs, &model.RegisterUserInvalidInputError{
				Message: errMessage.Error(),
				Field:   model.RegisterUserInvalidInputField(strings.ToLower(k)),
			})
		}

		return result, nil
	}

	result, err := u.RegisterUser(ctx, r.dbPool)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return result, nil
}

// LoginUser is the resolver for the loginUser field.
func (r *mutationResolver) LoginUser(ctx context.Context, input model.LoginUserInput) (model.LoginUserResult, error) {
	u := user.User{
		Email:    input.Email,
		Password: input.Password,
	}

	if err := u.AuthenticateUserValidate(); err != nil {
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)
		var result model.LoginUserInvalidInputErrors
		for k, errMessage := range errs {
			result.InvalidInputs = append(result.InvalidInputs, &model.LoginUserInvalidInputError{
				Message: errMessage.Error(),
				Field:   model.LoginUserInvalidInputField(strings.ToLower(k)),
			})
		}

		return result, nil
	}

	result, err := u.LoginUser(ctx, r.dbPool)
	if err != nil {
		return result, nil
	}
	return result, nil
}

// LogoutUser is the resolver for the logoutUser field.
func (r *mutationResolver) LogoutUser(ctx context.Context) (bool, error) {
	cookie.RemoveAuthCookie(ctx)
	return true, nil
}

// SendVerifyEmail is the resolver for the sendVerifyEmail field.
func (r *mutationResolver) SendVerifyEmail(ctx context.Context) (bool, error) {
	isSentEmail, err := user.SendVerifyEmail(ctx, r.dbPool)
	if err != nil {
		return false, err
	}
	return isSentEmail, nil
}

// SendVerifyNewEmail is the resolver for the sendVerifyNewEmail field.
func (r *mutationResolver) SendVerifyNewEmail(ctx context.Context, input model.SendVerifyNewEmailInput) (model.SendVerifyNewEmailResult, error) {
	u := user.User{
		Email: input.Email,
	}

	err := u.SendVerifyNewEmailValidate()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)

		var result model.SendVerifyNewEmailInvalidInputErrors
		for k, errMessage := range errs {
			result.InvalidInputs = append(result.InvalidInputs, &model.SendVerifyNewEmailInvalidInputError{
				Message: errMessage.Error(),
				Field:   model.SendVerifyNewEmailInvalidInputField(strings.ToLower(k)),
			})
		}
		return result, nil
	}

	result, err := u.SendVerifyNewEmail(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ChangePassword is the resolver for the changePassword field.
func (r *mutationResolver) ChangePassword(ctx context.Context, input model.ChangePasswordInput) (model.ChangePasswordResult, error) {
	i := user.ChangePasswordInput{
		CurrentPassword:         input.CurrentPassword,
		NewPassword:             input.NewPassword,
		NewPasswordConfirmation: input.NewPasswordConfirmation,
	}

	if err := i.ChangePasswordValidate(); err != nil {
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)

		var result model.ChangePasswordInvalidInputErrors
		for k, errMessage := range errs {
			result.InvalidInputs = append(result.InvalidInputs, &model.ChangePasswordInvalidInputError{
				Message: errMessage.Error(),
				Field:   model.ChangePasswordInvalidInputField(strings.ToLower(k[:1]) + k[1:]),
			})
		}
		return result, nil
	}

	result, err := i.ChangePassword(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// VerifyEmail is the resolver for the verifyEmail field.
func (r *mutationResolver) VerifyEmail(ctx context.Context, input model.VerifyEmailInput) (model.VerifyEmailResult, error) {
	i := user.VerifyEmailInput{
		Code: input.Code,
	}
	if err := i.VerifyEmailValidate(); err != nil {
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)

		var result model.VerifyEmailInvalidInputErrors
		for k, errMessage := range errs {
			result.InvalidInputs = append(result.InvalidInputs, &model.VerifyEmailInvalidInputError{
				Message: errMessage.Error(),
				Field:   model.VerifyEmailInvalidInputField(strings.ToLower(k)),
			})
		}
		return result, nil
	}

	result, err := i.VerifyEmail(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SendResetPasswordEmail is the resolver for the sendResetPasswordEmail field.
func (r *mutationResolver) SendResetPasswordEmail(ctx context.Context, input model.SendResetPasswordEmailInput) (model.SendResetPasswordEmailResult, error) {
	panic(fmt.Errorf("not implemented"))
}

// Viewer is the resolver for the viewer field.
func (r *queryResolver) Viewer(ctx context.Context) (*model.Viewer, error) {
	user := user.GetViewer(ctx)
	return &model.Viewer{AccountUser: user}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	user, err := user.GetUser(ctx, r.dbPool, utils.DecodeUniqueIDIdentifierOnly(id))
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ID is the resolver for the id field.
func (r *userResolver) ID(ctx context.Context, obj *model.User) (string, error) {
	return utils.GenerateUniqueID("User", obj.DatabaseID), nil
}

// EmailVerificationStatus is the resolver for the emailVerificationStatus field.
func (r *userResolver) EmailVerificationStatus(ctx context.Context, obj *model.User) (model.EmailVerificationStatus, error) {
	return model.EmailVerificationStatus(strings.ToUpper(obj.EmailVerificationStatus.String())), nil
}

// Recruitments is the resolver for the recruitments field.
func (r *userResolver) Recruitments(ctx context.Context, obj *model.User, first *int, after *string) (*model.RecruitmentConnection, error) {
	params, err := search.NewSearchParams(first, after, nil, nil)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	connection, err := recruitment.GetUserRecruitments(ctx, r.dbPool, obj.DatabaseID, params)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return connection, err
}

// Followings is the resolver for the followings field.
func (r *userResolver) Followings(ctx context.Context, obj *model.User, first *int, after *string) (*model.FollowConnection, error) {
	params, err := search.NewSearchParams(first, after, nil, nil)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	connection, err := relationship.GetFollowings(ctx, r.dbPool, obj.DatabaseID, params)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return connection, nil
}

// FeedbackFollow is the resolver for the feedbackFollow field.
func (r *userResolver) FeedbackFollow(ctx context.Context, obj *model.User) (*model.FeedbackFollow, error) {
	viewer := user.GetViewer(ctx)
	//* ログインしていない時とログインユーザーとobjが同じ場合は処理を進めない。follow,unfollowでviewer.accountUser.feedbackFollowを返すため
	if viewer == nil || viewer.DatabaseID == obj.DatabaseID {
		return &model.FeedbackFollow{UserID: obj.DatabaseID}, nil
	}
	feedback, err := relationship.GetFeedbackFollow(ctx, r.dbPool, obj.DatabaseID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return feedback, nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
