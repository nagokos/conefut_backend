package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/graph/cookie"
	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/prefecture"
	"github.com/nagokos/connefut_backend/graph/models/recruitment"
	"github.com/nagokos/connefut_backend/graph/models/relationship"
	"github.com/nagokos/connefut_backend/graph/models/search"
	"github.com/nagokos/connefut_backend/graph/models/sport"
	"github.com/nagokos/connefut_backend/graph/models/user"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

const ()

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

// ResendVerificationCodeToUser is the resolver for the resendVerificationCodeToUser field.
func (r *mutationResolver) ResendVerificationCodeToUser(ctx context.Context) (bool, error) {
	isSentEmail, err := user.SendVerifyEmail(ctx, r.dbPool)
	if err != nil {
		return false, err
	}
	return isSentEmail, nil
}

// ChangeUserEmail is the resolver for the changeUserEmail field.
func (r *mutationResolver) ChangeUserEmail(ctx context.Context, input model.ChangeUserEmailInput) (model.ChangeUserEmailResult, error) {
	i := user.ChangeEmailInput{
		NewEmail: input.NewEmail,
	}

	if err := i.ChangeEmailValidate(); err != nil {
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)

		for field, message := range errs {
			result := model.ChangeUserEmailInvalidInputError{
				Message: message.Error(),
				Field:   model.ChangeUserEmailInvalidInputField(strings.ToUpper(field)),
			}
			return result, nil
		}
	}

	result, err := i.ChangeEmail(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// VerifyUserEmail is the resolver for the verifyUserEmail field.
func (r *mutationResolver) VerifyUserEmail(ctx context.Context, input model.VerifyUserEmailInput) (model.VerifyUserEmailResult, error) {
	i := user.VerifyEmailInput{
		Code: input.Code,
	}
	if err := i.VerifyEmailValidate(); err != nil {
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)

		for field, message := range errs {
			result := model.VerifyUserEmailInvalidInputError{
				Message: message.Error(),
				Field:   model.VerifyUserEmailInvalidInputField(strings.ToUpper(field)),
			}
			return result, nil
		}
	}

	result, err := i.VerifyEmail(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ChangeUserPassword is the resolver for the changeUserPassword field.
func (r *mutationResolver) ChangeUserPassword(ctx context.Context, input model.ChangeUserPasswordInput) (model.ChangeUserPasswordResult, error) {
	i := user.ChangePasswordInput{
		CurrentPassword:         input.CurrentPassword,
		NewPassword:             input.NewPassword,
		NewPasswordConfirmation: input.NewPasswordConfirmation,
	}

	if err := i.ChangePasswordValidate(); err != nil {
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)

		var result model.ChangeUserPasswordInvalidInputErrors
		for k, errMessage := range errs {
			result.InvalidInputs = append(result.InvalidInputs, &model.ChangeUserPasswordInvalidInputError{
				Message: errMessage.Error(),
				Field:   model.ChangeUserPasswordInvalidInputField(strings.ToUpper(k[:1]) + k[1:]),
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

// SendResetPasswordEmailToUser is the resolver for the sendResetPasswordEmailToUser field.
func (r *mutationResolver) SendResetPasswordEmailToUser(ctx context.Context, input model.SendResetPasswordEmailToUserInput) (model.SendResetPasswordEmailToUserResult, error) {
	i := user.ResetPasswordInput{
		Email: input.Email,
	}
	if err := i.SendResetPasswordEmailValidate(); err != nil {
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)

		for field, message := range errs {
			result := model.SendResetPasswordEmailToUserInvalidInputError{
				Field:   model.SendResetPasswordEmailToUserInvalidInputField(strings.ToUpper(field)),
				Message: message.Error(),
			}
			return result, nil
		}
	}

	result, err := i.SendResetPasswordEmail(ctx, r.dbPool)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return result, nil
}

// ResetUserPassword is the resolver for the resetUserPassword field.
func (r *mutationResolver) ResetUserPassword(ctx context.Context, token string, input model.ResetUserPasswordInput) (model.ResetUserPasswordResult, error) {
	fmt.Println("h")
	i := user.ResetPasswordInput{
		NewPassword:             input.NewPassword,
		NewPasswordConfirmation: input.NewPasswordConfirmation,
	}

	if err := i.ResetPasswordValidate(); err != nil {
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)

		var result model.ResetUserPasswordInvalidInputErrors
		for field, message := range errs {
			result.InvalidInputs = append(result.InvalidInputs, &model.ResetUserPasswordInvalidInputError{
				Field:   model.ResetUserPasswordInvalidInputField(strings.ToUpper(field)),
				Message: message.Error(),
			})
		}
		return result, nil
	}

	result, err := i.ResetPassword(ctx, r.dbPool, token)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return result, nil
}

// UpdateUser is the resolver for the updateUser field.
func (r *mutationResolver) UpdateUser(ctx context.Context, input model.UpdateUserInput) (model.UpdateUserResult, error) {
	u := user.User{
		Name:          input.Name,
		Introduction:  input.Introduction,
		PrefectureIDs: utils.DecodeUniqueIDs(input.PrefectureIds),
		SportIDs:      utils.DecodeUniqueIDs(input.SportIds),
		WebsiteURL:    input.WebsiteURL,
	}

	if err := u.UpdateUserValidate(); err != nil {
		logger.NewLogger().Error(err.Error())
		errs := err.(validation.Errors)

		var result model.UpdateUserInvalidInputErrors
		for field, message := range errs {
			fmt.Println(strings.ToLower(field[:1]) + field[1:])
			result.InvalidInputs = append(result.InvalidInputs, &model.UpdateUserInvalidInputError{
				Field:   model.UpdateUserInvalidInputField(strings.ToLower(field[:1]) + field[1:]),
				Message: message.Error(),
			})
		}
		return result, nil
	}

	result, err := u.UpdateUser(ctx, r.dbPool)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return result, nil
}

// UploadUserAvatar is the resolver for the uploadUserAvatar field.
func (r *mutationResolver) UploadUserAvatar(ctx context.Context, input model.UploadUserAvatarInput) (*model.UploadUserAvatarSuccess, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	viewer := user.GetViewer(ctx)
	bucketName := "connefut-user-upload"
	objectPath := "avatar/"
	objectName := fmt.Sprintf("%s.%s", base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf("%d:%s", viewer.DatabaseID, input.File.Filename))), strings.Split(input.File.Filename, ".")[1])
	writer := client.Bucket(bucketName).Object(objectPath + objectName).NewWriter(ctx)
	if _, err := io.Copy(writer, input.File.File); err != nil {
		panic(err)
	}

	cmd := `
	  UPDATE users
		SET (avatar, updated_at) = ($1, $2)
		WHERE id = $3
		RETURNING id, avatar
	`
	publicPath := fmt.Sprintf("https://storage.googleapis.com/%s%s/%s", bucketName, writer.Bucket, writer.Name)
	now := time.Now().Local()
	row := r.dbPool.QueryRow(ctx, cmd, publicPath, now, viewer.DatabaseID)
	var user model.User
	if err := row.Scan(&user.DatabaseID, &user.Avatar); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	if publicPath != viewer.Avatar {
		o := client.Bucket(bucketName).Object(viewer.Avatar[strings.Index(viewer.Avatar, objectPath):])
		exists, _ := o.Attrs(ctx)
		if exists != nil {
			if err := o.Delete(ctx); err != nil {
				logger.NewLogger().Error(err.Error())
				return nil, err
			}
		}
	}

	return &model.UploadUserAvatarSuccess{Viewer: &model.Viewer{AccountUser: &user}}, nil
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

// PlaySports is the resolver for the playSports field.
func (r *userResolver) PlaySports(ctx context.Context, obj *model.User) ([]*model.Sport, error) {
	sports, err := sport.GetSportsByUserID(ctx, r.dbPool, obj.DatabaseID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return sports, nil
}

// ActivityAreas is the resolver for the activityAreas field.
func (r *userResolver) ActivityAreas(ctx context.Context, obj *model.User) ([]*model.Prefecture, error) {
	prefectures, err := prefecture.GetPrefecturesByUserID(ctx, r.dbPool, obj.DatabaseID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return prefectures, nil
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
