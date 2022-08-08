package relationship

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/search"
	"github.com/nagokos/connefut_backend/graph/models/user"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

func Follow(ctx context.Context, dbPool *pgxpool.Pool, userID int) (*model.FollowResult, error) {
	viewer := user.GetViewer(ctx)
	timeNow := time.Now().Local()
	cmd := `
	  INSERT INTO relationships 
		  (followed_id, follower_id, created_at, updated_at) 
		VALUES 
		  ($1, $2, $3, $4)
	`
	if _, err := dbPool.Exec(ctx, cmd, userID, viewer.DatabaseID, timeNow, timeNow); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	result := model.FollowResult{
		Viewer: &model.Viewer{
			AccountUser: viewer,
		},
		FeedbackFollow: &model.FeedbackFollow{
			UserID:           userID,
			IsViewerFollowed: true,
		},
	}

	return &result, nil
}

func UnFollow(ctx context.Context, dbPool *pgxpool.Pool, userID int) (*model.UnFollowResult, error) {
	viewer := user.GetViewer(ctx)
	cmd := `
	  DELETE FROM relationships
		WHERE followed_id = $1
		AND follower_id = $2
	`
	if _, err := dbPool.Exec(ctx, cmd, userID, viewer.DatabaseID); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	result := model.UnFollowResult{
		Viewer: &model.Viewer{
			AccountUser: viewer,
		},
		FeedbackFollow: &model.FeedbackFollow{
			UserID:           userID,
			IsViewerFollowed: false,
		},
	}
	return &result, nil
}

func GetFollowings(ctx context.Context, dbPool *pgxpool.Pool, userID int, params search.SearchParams) (*model.FollowConnection, error) {
	cmd := `
	  SELECT u.id, u.name, u.avatar
		FROM (
			SELECT followed_id
			FROM relationships
			WHERE follower_id = $1
			AND ($2 OR followed_id > $3)
			LIMIT $4
		) as r
		INNER JOIN users as u
		  ON u.id = r.followed_id
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

		connection.Edges = append(connection.Edges, &model.FollowEdge{
			Cursor: utils.GenerateUniqueID("User", following.DatabaseID),
			Node:   &following,
		})
	}

	if err := rows.Err(); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	if len(connection.Edges) > 0 {
		lastEdge := connection.Edges[len(connection.Edges)-1]

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
			userID, lastEdge.Node.DatabaseID,
		)

		var count int
		if err := row.Scan(&count); err != nil {
			logger.NewLogger().Error(err.Error())
			return nil, err
		}
		if count > 0 {
			connection.PageInfo.HasNextPage = true
		}
		connection.PageInfo.EndCursor = &lastEdge.Cursor
	}

	return &connection, nil
}

func GetFeedbackFollow(ctx context.Context, dbPool *pgxpool.Pool, userID int) (*model.FeedbackFollow, error) {
	isFollowed, err := IsViewerFollowed(ctx, dbPool, userID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	feedback := model.FeedbackFollow{
		UserID:           userID,
		IsViewerFollowed: isFollowed,
	}
	return &feedback, nil
}

func IsViewerFollowed(ctx context.Context, dbPool *pgxpool.Pool, userID int) (bool, error) {
	cmd := `
		SELECT COUNT(DISTINCT id)
		FROM relationships
		WHERE follower_id = $1
		AND followed_id = $2
	`
	viewer := user.GetViewer(ctx)
	row := dbPool.QueryRow(ctx, cmd, viewer.DatabaseID, userID)
	var count int
	if err := row.Scan(&count); err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func GetFollowingsCount(ctx context.Context, dbPool *pgxpool.Pool, userID int) (int, error) {
	cmd := `
	  SELECT COUNT(DISTINCT id)
		FROM relationships
		WHERE follower_id = $1
	`
	row := dbPool.QueryRow(ctx, cmd, userID)
	var count int
	if err := row.Scan(&count); err != nil {
		logger.NewLogger().Error(err.Error())
		return 0, nil
	}
	return count, nil
}
