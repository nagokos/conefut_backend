package applicant

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/rs/xid"
)

func CheckAppliedForRecruitment(ctx context.Context, dbPool *pgxpool.Pool, recID string) (bool, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return false, nil
	}

	cmd := `
	  SELECT COUNT(DISTINCT a.id)
		FROM applicants AS a 
		WHERE a.recruitment_id = $1
		AND a.user_id = $2
	`

	row := dbPool.QueryRow(ctx, cmd, recID, currentUser.ID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	var isApplied bool
	if count != 0 {
		isApplied = true
	}

	return isApplied, nil
}

func GetAppliedCounts(ctx context.Context, dbPool *pgxpool.Pool, recID string) (int, error) {
	cmd := `
	  SELECT COUNT(DISTINCT a.id)
		FROM applicants AS a
		WHERE a.recruitment_id = $1
	`
	row := dbPool.QueryRow(ctx, cmd, recID)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return count, err
	}

	return count, nil
}

func CreateApplicant(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID, message string) (bool, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return false, errors.New("ログインしてください")
	}

	cmd := "SELECT r.user_id FROM recruitments AS r WHERE r.id = $1"
	row := dbPool.QueryRow(ctx, cmd, recruitmentID)

	var userID string
	err := row.Scan(&userID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	if currentUser.ID == userID {
		logger.NewLogger().Error("This is a self-generated recruitment")
		return false, errors.New("自分が作成した募集には応募できません")
	}

	cmd = `
	  INSERT INTO applicants 
		  (id, recruitment_id, user_id, created_at, updated_at, message)
		VALUES 
		  ($1, $2, $3, $4, $5, $6)
	`

	timeNow := time.Now().Local()

	_, err = dbPool.Exec(
		ctx, cmd,
		xid.New().String(), recruitmentID, currentUser.ID, timeNow, timeNow, message,
	)

	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	return true, nil
}
