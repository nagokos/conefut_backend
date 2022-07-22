package relationship

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

func CheckFollowed(ctx context.Context, dbPool *pgxpool.Pool, userID string) (*model.FeedbackFollow, error) {
	feedback := model.FeedbackFollow{
		ID: utils.GenerateUniqueID("Relationship", utils.DecodeUniqueID(userID)),
	}

	cmd := `
	  SELECT COUNT(DISTINCT id)
		FROM relationships
		WHERE followed_id = $1
		AND follower_id = $2
	`
	viewer := auth.ForContext(ctx)
	row := dbPool.QueryRow(
		ctx, cmd,
		viewer.DatabaseID, utils.DecodeUniqueID(userID),
	)

	var count int
	err := row.Scan(&count)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	if count > 0 {
		feedback.ViewerDoesFollow = true
	}

	return &feedback, nil
}

func Follow(ctx context.Context, dbPool *pgxpool.Pool, userID string) (*model.FeedbackFollow, error) {
	feedback := model.FeedbackFollow{
		ID: utils.GenerateUniqueID("Relationship", utils.DecodeUniqueID(userID)),
	}

	viewer := auth.ForContext(ctx)
	timeNow := time.Now().Local()

	cmd := "INSERT INTO relationships (followed_id, follower_id, created_at, updated_at) VALUES ($1, $2, $3, $4)"
	_, err := dbPool.Exec(
		ctx, cmd,
		viewer.DatabaseID, utils.DecodeUniqueID(userID), timeNow, timeNow,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	feedback.ViewerDoesFollow = true

	return &feedback, nil
}

func UnFollow(ctx context.Context, dbPool *pgxpool.Pool, userID string) (*model.FeedbackFollow, error) {
	feedback := model.FeedbackFollow{
		ID: utils.GenerateUniqueID("Relationship", utils.DecodeUniqueID(userID)),
	}

	viewer := auth.ForContext(ctx)
	cmd := `
	  DELETE FROM relationships
		WHERE followed_id = $1
		AND follower_id = $2
	`
	_, err := dbPool.Exec(
		ctx, cmd,
		viewer.DatabaseID, utils.DecodeUniqueID(userID),
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &feedback, nil
}
