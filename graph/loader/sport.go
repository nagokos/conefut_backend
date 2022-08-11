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

type SportReader struct {
	dbPool *pgxpool.Pool
}

func (u *SportReader) GetSports(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	sportIDs := make([]interface{}, len(keys))
	var cmdArray []string
	for ix, key := range keys {
		sportIDs[ix] = key.String()
		cmdArray = append(cmdArray, fmt.Sprintf("$%d", ix+1))
	}
	cmd := fmt.Sprintf("SELECT id, name FROM sports WHERE id IN (%s)", strings.Join(cmdArray, ","))

	rows, err := u.dbPool.Query(
		ctx,
		cmd,
		sportIDs...,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil
	}
	defer rows.Close()

	sportByID := map[string]*model.Sport{}
	for rows.Next() {
		var sport model.Sport
		err := rows.Scan(&sport.DatabaseID, &sport.Name)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		sportByID[strconv.Itoa(sport.DatabaseID)] = &sport
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil
	}

	output := make([]*dataloader.Result, len(keys))
	for index, sportKey := range keys {
		sport, ok := sportByID[sportKey.String()]
		if ok {
			output[index] = &dataloader.Result{Data: sport, Error: nil}
		} else {
			err := fmt.Errorf("sport not found %s", sportKey.String())
			output[index] = &dataloader.Result{Data: nil, Error: err}
		}
	}
	return output
}

func GetSport(ctx context.Context, sportID int) (*model.Sport, error) {
	loaders := GetLoaders(ctx)
	thunk := loaders.SportLoader.Load(ctx, dataloader.StringKey(fmt.Sprintf("%d", sportID)))
	result, err := thunk()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return result.(*model.Sport), nil
}
