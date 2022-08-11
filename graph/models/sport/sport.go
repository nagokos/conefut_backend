package sport

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

func GetSports(ctx context.Context, dbPool *pgxpool.Pool) ([]*model.Sport, error) {
	cmd := "SELECT id, name FROM sports"
	rows, err := dbPool.Query(ctx, cmd)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var sports []*model.Sport
	for rows.Next() {
		var sport model.Sport
		err := rows.Scan(&sport.DatabaseID, &sport.Name)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		sports = append(sports, &sport)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return sports, nil
}

func GetSport(ctx context.Context, dbPool *pgxpool.Pool, id int) (*model.Sport, error) {
	cmd := "SELECT id, name FROM sports WHERE id = $1"

	var sport model.Sport
	row := dbPool.QueryRow(ctx, cmd, id)
	err := row.Scan(&sport.DatabaseID, &sport.Name)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return &sport, nil
}
