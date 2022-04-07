package competition

import (
	"database/sql"

	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

func GetCompetitions(dbConnection *sql.DB) ([]*model.Competition, error) {
	cmd := "SELECT id, name FROM competitions"
	rows, err := dbConnection.Query(cmd)
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
