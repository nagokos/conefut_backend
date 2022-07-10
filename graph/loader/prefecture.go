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

type PrefectureReader struct {
	dbPool *pgxpool.Pool
}

func (u *PrefectureReader) GetPrefectures(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	prefectureIDs := make([]interface{}, len(keys))
	var cmdArray []string
	for ix, key := range keys {
		prefectureIDs[ix] = key.String()
		cmdArray = append(cmdArray, fmt.Sprintf("$%d", ix+1))
	}
	cmd := fmt.Sprintf("SELECT id, name FROM prefectures WHERE id IN (%s)", strings.Join(cmdArray, ","))

	rows, err := u.dbPool.Query(
		ctx,
		cmd,
		prefectureIDs...,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil
	}
	defer rows.Close()

	prefectureByID := map[string]*model.Prefecture{}
	for rows.Next() {
		var prefecture model.Prefecture
		err := rows.Scan(&prefecture.DatabaseID, &prefecture.Name)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		prefectureByID[strconv.Itoa(prefecture.DatabaseID)] = &prefecture
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil
	}

	output := make([]*dataloader.Result, len(keys))
	for index, prefectureKey := range keys {
		prefecture, ok := prefectureByID[prefectureKey.String()]
		if ok {
			output[index] = &dataloader.Result{Data: prefecture, Error: nil}
		} else {
			err := fmt.Errorf("prefecture not found %s", prefectureKey.String())
			output[index] = &dataloader.Result{Data: nil, Error: err}
		}
	}
	return output
}

func GetPrefecture(ctx context.Context, prefectureID int) (*model.Prefecture, error) {
	loaders := GetLoaders(ctx)
	thunk := loaders.PrefectureLoader.Load(ctx, dataloader.StringKey(fmt.Sprintf("%d", prefectureID)))
	result, err := thunk()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return result.(*model.Prefecture), nil
}
