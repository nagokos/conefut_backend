package prefecture

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

type NullablePrefecture struct {
	ID         *string
	DatabaseID *int
	Name       *string
}

func GetPrefectures(ctx context.Context, dbPool *pgxpool.Pool) ([]*model.Prefecture, error) {
	var prefectures []*model.Prefecture

	cmd := "SELECT id, name FROM prefectures"
	rows, err := dbPool.Query(ctx, cmd)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var prefecture model.Prefecture
		err := rows.Scan(&prefecture.DatabaseID, &prefecture.Name)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		prefectures = append(prefectures, &prefecture)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return prefectures, nil
}

func GetPrefecture(ctx context.Context, dbPool *pgxpool.Pool, id int) (*model.Prefecture, error) {
	cmd := "SELECT id, name FROM prefectures WHERE id = $1"

	var prefecture model.Prefecture
	row := dbPool.QueryRow(ctx, cmd, id)
	err := row.Scan(&prefecture.DatabaseID, &prefecture.Name)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &prefecture, nil
}
