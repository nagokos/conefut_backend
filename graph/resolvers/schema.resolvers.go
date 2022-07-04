package resolvers

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
	"github.com/nagokos/connefut_backend/graph/models/entrie"
	"github.com/nagokos/connefut_backend/graph/models/message"
	"github.com/nagokos/connefut_backend/graph/models/room"
	"github.com/nagokos/connefut_backend/graph/models/stock"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

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

func (r *queryResolver) Node(ctx context.Context, id string) (model.Node, error) {
	panic(fmt.Errorf("not implemented"))
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
