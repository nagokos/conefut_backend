package applicant

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/room"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/rs/xid"
)

func GetApplicant(ctx context.Context, dbPool *pgxpool.Pool, id int) (*model.Applicant, error) {
	cmd := "SELECT id, message FROM applicants WHERE id = $1"

	var applicant model.Applicant
	row := dbPool.QueryRow(ctx, cmd, id)
	err := row.Scan(&applicant.DatabaseID, &applicant.Message)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &applicant, nil
}

func CheckAppliedForRecruitment(ctx context.Context, dbPool *pgxpool.Pool, recID string) (bool, error) {
	viewer := auth.ForContext(ctx)
	if viewer == nil {
		return false, nil
	}

	cmd := `
	  SELECT COUNT(DISTINCT a.id)
		FROM applicants AS a 
		WHERE a.recruitment_id = $1
		AND a.user_id = $2
	`

	row := dbPool.QueryRow(ctx, cmd, recID, viewer.ID)

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

func CreateApplicant(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID, message string) (bool, error) {
	viewer := auth.ForContext(ctx)
	if viewer == nil {
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

	if viewer.ID == userID {
		logger.NewLogger().Error("This is a self-generated recruitment")
		return false, errors.New("自分が作成した募集には応募できません")
	}

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	cmd = `
	  INSERT INTO applicants 
		  (id, recruitment_id, user_id, created_at, updated_at, message)
		VALUES 
		  ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	timeNow := time.Now().Local()

	row = tx.QueryRow(
		ctx, cmd,
		xid.New().String(), recruitmentID, viewer.ID, timeNow, timeNow, message,
	)

	var applicantID string
	err = row.Scan(&applicantID)

	if err != nil {
		logger.NewLogger().Error(err.Error())
		err = tx.Rollback(ctx)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		return false, err
	}

	cmd = `
	  SELECT e_1.room_id
		FROM entries AS e_1
		WHERE e_1.user_id = $1
		AND EXISTS(
		  SELECT 1
			FROM entries AS e_2
			WHERE e_1.room_id = e_2.room_id
			AND e_2.user_id = $2
		)
	`
	row = tx.QueryRow(
		ctx, cmd,
		viewer.ID, userID,
	)

	var roomID string
	err = row.Scan(&roomID)

	if err != nil {
		if err == pgx.ErrNoRows {
			roomID, err = room.CreateRoom(ctx, tx)
			if err != nil {
				err = tx.Rollback(ctx)
				if err != nil {
					logger.NewLogger().Error(err.Error())
				}
				return false, err
			}

			entrieUsers := [2]string{viewer.ID, userID}
			cmd = `
				INSERT INTO entries
					(id, room_id, user_id, created_at, updated_at)
				VALUES
					($1, $2, $3, $4, $5)
			`

			for _, userID := range entrieUsers {
				_, err = tx.Exec(
					ctx, cmd,
					xid.New().String(), roomID, userID, timeNow, timeNow,
				)
				if err != nil {
					logger.NewLogger().Error(err.Error())
					err = tx.Rollback(ctx)
					if err != nil {
						logger.NewLogger().Error(err.Error())
					}
					return false, err
				}
			}
		} else {
			logger.NewLogger().Error(err.Error())
			err = tx.Rollback(ctx)
			if err != nil {
				logger.NewLogger().Error(err.Error())
			}
			return false, err
		}
	}

	cmd = `
	  INSERT INTO messages
		  (id, room_id, user_id, applicant_id, created_at, updated_at)
		VALUES
		  ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(
		ctx, cmd,
		xid.New().String(), roomID, viewer.ID, applicantID, timeNow, timeNow,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		err = tx.Rollback(ctx)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		return false, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		err = tx.Rollback(ctx)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		return false, err
	}

	return true, nil
}
