package applicant

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

func CheckApplied(ctx context.Context, dbConnection *sql.DB, recID string) (bool, error) {
	currentUser := auth.ForContext(ctx)

	cmd := `
	  SELECT COUNT(DISTINCT a.id) 
		FROM applicants AS a
		WHERE a.user_id = $1
		AND a.recruitment_id = $2
	`

	row := dbConnection.QueryRow(cmd, currentUser.ID, recID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	var isApplied bool
	if count == 1 {
		isApplied = true
	}

	return isApplied, nil
}

func GetAppliedCounts(dbConnection *sql.DB, recID string) (int, error) {
	cmd := `
	  SELECT COUNT(DISTINCT a.id)
		FROM applicants AS a
		WHERE a.recruitment_id = $1
	`
	row := dbConnection.QueryRow(cmd, recID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return count, err
	}

	return count, nil
}

func CreateApplicant(ctx context.Context, dbConnection *sql.DB, recId string, status model.ManagementStatus) (bool, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return false, errors.New("ログインしてください")
	}

	return true, nil
}
