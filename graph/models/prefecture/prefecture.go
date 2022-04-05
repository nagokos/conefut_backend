package prefecture

import (
	"database/sql"

	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

func GetPrefectures(dbConnection *sql.DB) ([]*model.Prefecture, error) {
	var prefectures []*model.Prefecture

	cmd := "SELECT id, name FROM prefectures"
	rows, err := dbConnection.Query(cmd)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var prefecture model.Prefecture
		err := rows.Scan(&prefecture.ID, &prefecture.Name)
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
