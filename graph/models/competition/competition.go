package competition

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

type NullableCompetition struct {
	ID   *string
	Name *string
}

func GetCompetitions(ctx context.Context, dbPool *pgxpool.Pool) ([]*model.Competition, error) {
	cmd := "SELECT id, name FROM competitions"
	rows, err := dbPool.Query(ctx, cmd)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var competitions []*model.Competition
	for rows.Next() {
		var competition model.Competition
		err := rows.Scan(&competition.ID, &competition.Name)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		competitions = append(competitions, &competition)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return competitions, nil
}
