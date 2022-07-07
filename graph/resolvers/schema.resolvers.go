package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/graph/generated"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/entrie"
	"github.com/nagokos/connefut_backend/graph/models/message"
	"github.com/nagokos/connefut_backend/graph/models/room"
	"github.com/nagokos/connefut_backend/graph/models/stock"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

// CreateStock is the resolver for the createStock field.
func (r *mutationResolver) CreateStock(ctx context.Context, recruitmentID string) (bool, error) {
	_, err := stock.CreateStock(ctx, r.dbPool, recruitmentID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteStock is the resolver for the deleteStock field.
func (r *mutationResolver) DeleteStock(ctx context.Context, recruitmentID string) (bool, error) {
	_, err := stock.DeleteStock(ctx, r.dbPool, recruitmentID)
	if err != nil {
		return false, err
	}
	return true, nil
}

// AddRecruitmentTag is the resolver for the addRecruitmentTag field.
func (r *mutationResolver) AddRecruitmentTag(ctx context.Context, tagID string, recruitmentID string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

// CreateMessage is the resolver for the createMessage field.
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

// Node is the resolver for the node field.
func (r *queryResolver) Node(ctx context.Context, id string) (model.Node, error) {
	panic(fmt.Errorf("not implemented"))
}

// CheckStocked is the resolver for the checkStocked field.
func (r *queryResolver) CheckStocked(ctx context.Context, recruitmentID string) (bool, error) {
	res, err := stock.CheckStocked(ctx, r.dbPool, recruitmentID)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetStockedCount is the resolver for the getStockedCount field.
func (r *queryResolver) GetStockedCount(ctx context.Context, recruitmentID string) (int, error) {
	res, err := stock.GetStockedCount(ctx, r.dbPool, recruitmentID)
	if err != nil {
		return 0, err
	}
	return res, nil
}

// GetCurrentUserRooms is the resolver for the getCurrentUserRooms field.
func (r *queryResolver) GetCurrentUserRooms(ctx context.Context) ([]*model.Room, error) {
	res, err := room.GetCurrentUserRooms(ctx, r.dbPool)
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetEntrieUser is the resolver for the getEntrieUser field.
func (r *queryResolver) GetEntrieUser(ctx context.Context, roomID string) (*model.User, error) {
	res, err := entrie.GetEntrieUser(ctx, r.dbPool, roomID)
	if err != nil {
		return res, err
	}

	return res, nil
}

// GetRoomMessages is the resolver for the getRoomMessages field.
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
