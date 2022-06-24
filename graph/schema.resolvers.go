package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/applicant"
	"github.com/nagokos/connefut_backend/graph/models/competition"
	"github.com/nagokos/connefut_backend/graph/models/entrie"
	"github.com/nagokos/connefut_backend/graph/models/message"
	"github.com/nagokos/connefut_backend/graph/models/prefecture"
	"github.com/nagokos/connefut_backend/graph/models/recruitment"
	"github.com/nagokos/connefut_backend/graph/models/room"
	"github.com/nagokos/connefut_backend/graph/models/search"
	"github.com/nagokos/connefut_backend/graph/models/stock"
	"github.com/nagokos/connefut_backend/graph/models/tag"
	"github.com/nagokos/connefut_backend/graph/models/user"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.CreateUserInput) (bool, error) {
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

		return false, errors.New("フォームに不備があります")
	}

	userID, err := u.Insert(ctx, r.dbPool)
	if err != nil {
		return false, err
	}

	token, _ := user.CreateToken(userID)

	auth.SetAuthCookie(ctx, token)

	return true, nil
}

func (r *mutationResolver) LoginUser(ctx context.Context, input model.LoginUserInput) (bool, error) {
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

		return false, errors.New("フォームに不備があります")
	}

	userID, err := u.Authenticate(ctx, r.dbPool)
	if err != nil {
		return false, err
	}

	token, _ := user.CreateToken(userID)

	auth.SetAuthCookie(ctx, token)

	return true, nil
}

func (r *mutationResolver) LogoutUser(ctx context.Context) (bool, error) {
	auth.RemoveAuthCookie(ctx)
	return true, nil
}

func (r *mutationResolver) CreateRecruitment(ctx context.Context, input model.RecruitmentInput) (*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)
	if model.Status(input.Status) == model.StatusPublished &&
		currentUser.EmailVerificationStatus == model.EmailVerificationStatusPending {
		return &model.Recruitment{}, errors.New("メールアドレスを認証してください")
	}

	rm := recruitment.Recruitment{
		Title:         input.Title,
		Type:          input.Type,
		Detail:        input.Detail,
		StartAt:       input.StartAt,
		Place:         input.Place,
		LocationLat:   input.LocationLat,
		LocationLng:   input.LocationLng,
		Status:        input.Status,
		ClosingAt:     input.ClosingAt,
		CompetitionID: input.CompetitionID,
		PrefectureID:  input.PrefectureID,
		Tags:          input.Tags,
	}

	err := rm.RecruitmentValidate()
	if err != nil {
		logger.NewLogger().Sugar().Errorf("recruitment validation errors:", err.Error())
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(errMessage.Error(), utils.WithField(strings.ToLower(k))).AddGraphQLError(ctx)
		}
		return &model.Recruitment{}, err
	}

	resRecruitment, err := rm.CreateRecruitment(ctx, r.dbPool)
	if err != nil {
		return &model.Recruitment{}, err
	}

	return resRecruitment, nil
}

func (r *mutationResolver) UpdateRecruitment(ctx context.Context, id string, input model.RecruitmentInput) (*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)
	if model.Status(input.Status) == model.StatusPublished &&
		currentUser.EmailVerificationStatus == model.EmailVerificationStatusPending {
		return &model.Recruitment{}, errors.New("メールアドレスを認証してください")
	}

	rm := recruitment.Recruitment{
		Title:         input.Title,
		Type:          input.Type,
		Detail:        input.Detail,
		StartAt:       input.StartAt,
		Place:         input.Place,
		LocationLat:   input.LocationLat,
		LocationLng:   input.LocationLng,
		Status:        input.Status,
		ClosingAt:     input.ClosingAt,
		CompetitionID: input.CompetitionID,
		PrefectureID:  input.PrefectureID,
		Tags:          input.Tags,
	}

	err := rm.RecruitmentValidate()
	if err != nil {
		logger.NewLogger().Sugar().Errorf("recruitment validation errors %s", err.Error())
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(errMessage.Error(), utils.WithField(strings.ToLower(k))).AddGraphQLError(ctx)
		}
		return &model.Recruitment{}, err
	}

	res, err := rm.UpdateRecruitment(ctx, r.dbPool, id)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return res, err
	}

	return res, nil
}

func (r *mutationResolver) DeleteRecruitment(ctx context.Context, id string) (*model.Recruitment, error) {
	res, err := recruitment.DeleteRecruitment(ctx, r.dbPool, id)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (r *mutationResolver) CreateStock(ctx context.Context, recruitmentID string) (bool, error) {
	_, err := stock.CreateStock(ctx, r.dbPool, recruitmentID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) DeleteStock(ctx context.Context, recruitmentID string) (bool, error) {
	_, err := stock.DeleteStock(ctx, r.dbPool, recruitmentID)
	if err != nil {
		return false, err
	}
	return true, nil
}

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

func (r *mutationResolver) AddRecruitmentTag(ctx context.Context, tagID string, recruitmentID string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ApplyForRecruitment(ctx context.Context, recruitmentID string, input *model.ApplicantInput) (bool, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser.EmailVerificationStatus == model.EmailVerificationStatusPending {
		logger.NewLogger().Error("not email verified")
		return false, errors.New("メールアドレスを認証してください")
	}

	res, err := applicant.CreateApplicant(ctx, r.dbPool, recruitmentID, input.Message)
	if err != nil {
		return res, err
	}

	return res, err
}

func (r *mutationResolver) CreateMessage(ctx context.Context, roomID string, input model.CreateMessageInput) (*model.Message, error) {
	m := message.Message{
		Content: input.Content,
	}

	err := m.MessageValidate()
	if err != nil {
		logger.NewLogger().Sugar().Errorf("recruitment validation errors %s", err.Error())
		errs := err.(validation.Errors)

		for k, errMessage := range errs {
			utils.NewValidationError(errMessage.Error(), utils.WithField(strings.ToLower(k))).AddGraphQLError(ctx)
		}
		return nil, err
	}

	res, err := m.CreateMessage(ctx, r.dbPool, roomID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return res, nil
}

func (r *queryResolver) GetPrefectures(ctx context.Context) ([]*model.Prefecture, error) {
	res, err := prefecture.GetPrefectures(ctx, r.dbPool)
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

func (r *queryResolver) GetCompetitions(ctx context.Context) ([]*model.Competition, error) {
	res, err := competition.GetCompetitions(ctx, r.dbPool)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *queryResolver) GetRecruitments(ctx context.Context, input model.PaginationInput) (*model.RecruitmentConnection, error) {
	sp, err := search.NewSearchParams(input.After, input.Before, input.First, input.Last, input.Options)
	if err != nil {
		return nil, err
	}

	res, err := recruitment.GetRecruitments(ctx, r.dbPool, sp)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *queryResolver) GetCurrentUserRecruitments(ctx context.Context) ([]*model.Recruitment, error) {
	res, err := recruitment.GetCurrentUserRecruitments(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *queryResolver) GetRecruitment(ctx context.Context, id string) (*model.Recruitment, error) {
	res, err := recruitment.GetRecruitment(ctx, r.dbPool, id)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (r *queryResolver) GetStockedRecruitments(ctx context.Context) ([]*model.Recruitment, error) {
	res, err := recruitment.GetStockedRecruitments(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (r *queryResolver) GetAppliedRecruitments(ctx context.Context) ([]*model.Recruitment, error) {
	res, err := recruitment.GetAppliedRecruitments(ctx, r.dbPool)
	if err != nil {
		return res, err
	}

	return res, err
}

func (r *queryResolver) CheckStocked(ctx context.Context, recruitmentID string) (bool, error) {
	res, err := stock.CheckStocked(ctx, r.dbPool, recruitmentID)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (r *queryResolver) GetStockedCount(ctx context.Context, recruitmentID string) (int, error) {
	res, err := stock.GetStockedCount(ctx, r.dbPool, recruitmentID)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *queryResolver) GetTags(ctx context.Context) ([]*model.Tag, error) {
	res, err := tag.GetTags(ctx, r.dbPool)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *queryResolver) CheckAppliedForRecruitment(ctx context.Context, recruitmentID string) (bool, error) {
	res, err := applicant.CheckAppliedForRecruitment(ctx, r.dbPool, recruitmentID)
	if err != nil {
		return res, err
	}
	return res, err
}

func (r *queryResolver) GetAppliedCounts(ctx context.Context, recruitmentID string) (int, error) {
	res, err := applicant.GetAppliedCounts(ctx, r.dbPool, recruitmentID)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *queryResolver) GetCurrentUserRooms(ctx context.Context) ([]*model.Room, error) {
	res, err := room.GetCurrentUserRooms(ctx, r.dbPool)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *queryResolver) GetEntrieUser(ctx context.Context, roomID string) (*model.User, error) {
	res, err := entrie.GetEntrieUser(ctx, r.dbPool, roomID)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (r *queryResolver) GetRoomMessages(ctx context.Context, roomID string) ([]*model.Message, error) {
	res, err := message.GetRoomMessages(ctx, r.dbPool, roomID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return res, err
	}

	return res, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
