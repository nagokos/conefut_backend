package relationship

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/search"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

func CheckFollowed(ctx context.Context, dbPool *pgxpool.Pool, userID string) (*model.FeedbackFollow, error) {
	feedback := model.FeedbackFollow{
		ID: utils.GenerateUniqueID("Relationship", utils.DecodeUniqueID(userID)),
	}

	viewer := auth.ForContext(ctx)
	if viewer == nil {
		return &feedback, nil
	}

	cmd := `
	  SELECT COUNT(DISTINCT id)
		FROM relationships
		WHERE followed_id = $1
		AND follower_id = $2
	`
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

func CheckFollowedByRecruitmentID(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID string) (*model.FeedbackFollow, error) {
	cmd := `
	  SELECT id
		FROM users 
		WHERE id = (
			SELECT user_id
			FROM recruitments 
			WHERE id = $1
		)
	`
	row := dbPool.QueryRow(
		ctx, cmd,
		utils.DecodeUniqueID(recruitmentID),
	)

	var userID int
	err := row.Scan(&userID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	feedback := model.FeedbackFollow{
		ID: utils.GenerateUniqueID("Relationship", userID),
	}
	viewer := auth.ForContext(ctx)
	if viewer == nil {
		return &feedback, nil
	}

	cmd = `
	  SELECT COUNT(DISTINCT id)
		FROM relationships
		WHERE followed_id = $1
		AND follower_id = (
			SELECT user_id
			FROM recruitments 
			WHERE id = $2
		)
	`
	row = dbPool.QueryRow(
		ctx, cmd,
		viewer.DatabaseID, utils.DecodeUniqueID(recruitmentID),
	)

	var count int
	err = row.Scan(&count)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	if count > 0 {
		feedback.ViewerDoesFollow = true
	}

	feedback.ID = utils.GenerateUniqueID("Relationship", userID)
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

func GetFollowings(ctx context.Context, dbPool *pgxpool.Pool, userID int, params search.SearchParams) (*model.FollowConnection, error) {
	cmd := `
	  SELECT u.id, u.name, u.avatar
		FROM (
			SELECT id, follower_id
			FROM relationships
			WHERE followed_id = $1
			AND ($2 OR follower_id > $3)
			LIMIT $4
		) as r
		INNER JOIN users as u
		  ON u.id = r.follower_id
	`

	rows, err := dbPool.Query(
		ctx, cmd,
		userID, !params.UseAfter, params.After, params.NumRows,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	defer rows.Close()

	connection := model.FollowConnection{
		PageInfo: &model.PageInfo{},
	}
	for rows.Next() {
		var following model.User
		if err := rows.Scan(&following.DatabaseID, &following.Name, &following.Avatar); err != nil {
			logger.NewLogger().Error(err.Error())
		}

		feedback, err := CheckFollowed(ctx, dbPool, utils.GenerateUniqueID("User", following.DatabaseID))
		if err != nil {
			logger.NewLogger().Error(err.Error())
			return nil, err
		}

		connection.Edges = append(connection.Edges, &model.FollowEdge{
			Cursor:   utils.GenerateUniqueID("User", following.DatabaseID),
			Node:     &following,
			Feedback: feedback,
		})
	}

	if err := rows.Err(); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	cmd = `
	  SELECT COUNT(DISTINCT r.follower_id)
		FROM (
			SELECT follower_id
			FROM relationships
			WHERE followed_id = $1
		) as r
	`
	row := dbPool.QueryRow(ctx, cmd, userID)

	var totalCount int
	if err := row.Scan(&totalCount); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	connection.FollowCount = totalCount

	if len(connection.Edges) > 0 {
		endCursor := connection.Edges[len(connection.Edges)-1].Cursor
		connection.PageInfo.EndCursor = &endCursor

		cmd = `
			SELECT COUNT(DISTINCT r.follower_id)
			FROM (
				SELECT follower_id
				FROM relationships
				WHERE followed_id = $1
				AND follower_id > $2
				LIMIT 1
			) as r
		`
		row := dbPool.QueryRow(
			ctx, cmd,
			userID, utils.DecodeUniqueID(endCursor),
		)

		var count int
		if err := row.Scan(&count); err != nil {
			logger.NewLogger().Error(err.Error())
			return nil, err
		}

		if count > 0 {
			connection.PageInfo.HasNextPage = true
		}
	}

	return &connection, nil
}
