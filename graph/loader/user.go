package loader

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/graph-gophers/dataloader"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

type UserReader struct {
	dbPool *pgxpool.Pool
}

func (u *UserReader) GetUsers(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	userIDs := make([]interface{}, len(keys))
	var cmdArray []string
	for ix, key := range keys {
		userIDs[ix] = key.String()
		cmdArray = append(cmdArray, fmt.Sprintf("$%d", ix+1))
	}
	cmd := fmt.Sprintf("SELECT id, name, avatar FROM users WHERE id IN (%s)", strings.Join(cmdArray, ","))

	rows, err := u.dbPool.Query(
		ctx,
		cmd,
		userIDs...,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil
	}
	defer rows.Close()

	userByID := map[string]*model.User{}
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.DatabaseID, &user.Name, &user.Avatar)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		userByID[strconv.Itoa(user.DatabaseID)] = &user
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil
	}

	output := make([]*dataloader.Result, len(keys))
	for index, userKey := range keys {
		user, ok := userByID[userKey.String()]
		if ok {
			output[index] = &dataloader.Result{Data: user, Error: nil}
		} else {
			err := fmt.Errorf("user not found %s", userKey.String())
			output[index] = &dataloader.Result{Data: nil, Error: err}
		}
	}
	return output
}

func GetUser(ctx context.Context, userID int) (*model.User, error) {
	loaders := GetLoaders(ctx)
	thunk := loaders.UserLoader.Load(ctx, dataloader.StringKey(fmt.Sprintf("%d", userID)))
	result, err := thunk()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return result.(*model.User), nil
}
