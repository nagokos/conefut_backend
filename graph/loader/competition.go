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

type CompetitionReader struct {
	dbPool *pgxpool.Pool
}

func (u *CompetitionReader) GetCompetitions(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	competitionIDs := make([]interface{}, len(keys))
	var cmdArray []string
	for ix, key := range keys {
		competitionIDs[ix] = key.String()
		cmdArray = append(cmdArray, fmt.Sprintf("$%d", ix+1))
	}
	cmd := fmt.Sprintf("SELECT id, name FROM competitions WHERE id IN (%s)", strings.Join(cmdArray, ","))

	rows, err := u.dbPool.Query(
		ctx,
		cmd,
		competitionIDs...,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil
	}
	defer rows.Close()

	competitionByID := map[string]*model.Competition{}
	for rows.Next() {
		var competition model.Competition
		err := rows.Scan(&competition.DatabaseID, &competition.Name)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		competitionByID[strconv.Itoa(competition.DatabaseID)] = &competition
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil
	}

	output := make([]*dataloader.Result, len(keys))
	for index, competitionKey := range keys {
		competition, ok := competitionByID[competitionKey.String()]
		if ok {
			output[index] = &dataloader.Result{Data: competition, Error: nil}
		} else {
			err := fmt.Errorf("competition not found %s", competitionKey.String())
			output[index] = &dataloader.Result{Data: nil, Error: err}
		}
	}
	return output
}

func GetCompetition(ctx context.Context, competitionID int) (*model.Competition, error) {
	loaders := GetLoaders(ctx)
	thunk := loaders.CompetitionLoader.Load(ctx, dataloader.StringKey(fmt.Sprintf("%d", competitionID)))
	result, err := thunk()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return result.(*model.Competition), nil
}
