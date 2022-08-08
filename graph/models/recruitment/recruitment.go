package recruitment

import (
	"context"
	"errors"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/search"
	"github.com/nagokos/connefut_backend/graph/models/tag"
	"github.com/nagokos/connefut_backend/graph/models/user"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
)

type RecruitmentInput struct {
	Title         string
	Type          model.Type
	Venue         *string
	StartAt       *time.Time
	Detail        *string
	LocationLat   *float64
	LocationLng   *float64
	Status        model.Status
	ClosingAt     *time.Time
	CompetitionID int
	PrefectureID  int
	TagIDs        []int
}

func checkWithinTheDeadline(start time.Time) validation.RuleFunc {
	return func(v interface{}) error {
		var err error
		switch s := v.(type) {
		case *time.Time:
			difference := start.Sub(*s)
			if difference < 0 {
				err = errors.New("募集期限は開催日時よりも前に設定してください")
			}
		}
		return err
	}
}

func beforeNowStart(v interface{}) error {
	var err error
	switch t := v.(type) {
	case *time.Time:
		difference := time.Since(*t).Minutes()
		if difference >= 1 {
			err = errors.New("開催日時は現在以降に設定してください")
		} else {
			err = nil
		}
	}
	return err
}

func beforeNowClosing(v interface{}) error {
	var err error
	switch t := v.(type) {
	case *time.Time:
		difference := time.Since(*t).Minutes()
		if difference >= 1 {
			err = errors.New("募集期限は現在以降に設定してください")
		} else {
			err = nil
		}
	}
	return err
}

func (i RecruitmentInput) RecruitmentValidate() error {
	return validation.ValidateStruct(&i,
		validation.Field(
			&i.Title,
			validation.Required.Error("タイトルを入力してください"),
			validation.RuneLength(1, 60).Error("タイトルは1文字以上60文字以内で入力してください"),
		),
		validation.Field(
			&i.CompetitionID,
			validation.When(i.Status == model.StatusPublished,
				validation.Required.Error("募集競技を選択してください"),
			),
		),
		validation.Field(
			&i.Type,
			validation.In(
				model.TypeOpponent,
				model.TypePersonal,
				model.TypeMember,
				model.TypeJoin,
				model.TypeOther,
			),
		),
		validation.Field(
			&i.Detail,
			validation.When(i.Status == model.StatusPublished,
				validation.Required.Error("募集の詳細を入力してください"),
				validation.RuneLength(1, 10000).Error("募集の詳細は10000文字以内で入力してください"),
			).Else(validation.RuneLength(0, 10000).Error("募集の詳細は10000文字以内で入力してください")),
		),
		validation.Field(
			&i.PrefectureID,
			validation.Required.Error("募集エリアを選択してください"),
		),
		validation.Field(
			&i.Venue,
			validation.When(i.Status == model.StatusPublished,
				validation.When(i.Type == model.TypeOpponent || i.Type == model.TypePersonal,
					validation.Required.Error("会場名を入力してください"),
				),
			),
		),
		validation.Field(
			&i.StartAt,
			validation.When(i.Status == model.StatusPublished,
				validation.When(i.Type == model.TypeOpponent || i.Type == model.TypePersonal,
					validation.By(beforeNowStart),
					validation.Required.Error("開催日時を設定してください"),
				),
			),
		),
		validation.Field(
			&i.ClosingAt,
			validation.When(i.Status == model.StatusPublished,
				validation.Required.Error("募集期限を設定してください"),
				validation.When(i.Type == model.TypeOpponent || i.Type == model.TypePersonal,
					validation.By(beforeNowClosing),
					validation.By(checkWithinTheDeadline(*i.StartAt)),
				),
			),
		),
	)
}

//* 募集の作成
func (i *RecruitmentInput) CreateRecruitment(ctx context.Context, dbPool *pgxpool.Pool) (model.CreateRecruitmentResult, error) {
	cmd := `
	  INSERT INTO recruitments 
		  (title, competition_id, type, detail, prefecture_id, venue, start_at, closing_at, location_lat, location_lng, status, user_id, created_at, updated_at, published_at)
		VALUES
		  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, title, status, created_at, published_at, competition_id, user_id, closing_at, type, prefecture_id
		`
	timeNow := time.Now().Local()
	viewer := user.GetViewer(ctx)
	var publishedAt *time.Time
	if i.Status == model.StatusPublished {
		publishedAt = &timeNow
	} else {
		publishedAt = nil
	}
	row := dbPool.QueryRow(
		ctx, cmd,
		i.Title, i.CompetitionID, strings.ToLower(string(i.Type)), i.Detail, i.PrefectureID, i.Venue, i.StartAt,
		i.ClosingAt, i.LocationLat, i.LocationLng, strings.ToLower(string(i.Status)), viewer.DatabaseID, timeNow, timeNow, publishedAt,
	)

	var recruitment model.Recruitment
	if err := row.Scan(&recruitment.DatabaseID, &recruitment.Title, &recruitment.Status, &recruitment.CreatedAt, &recruitment.PublishedAt,
		&recruitment.CompetitionID, &recruitment.UserID, &recruitment.ClosingAt, &recruitment.Type, &recruitment.PrefectureID,
	); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	cmd = `
		INSERT INTO recruitment_tags
			(recruitment_id, tag_id, created_at, updated_at)
		VALUES
			($1, $2, $3, $4)
  `
	if len(i.TagIDs) > 0 {
		for _, sentTag := range i.TagIDs {
			if _, err := dbPool.Exec(ctx, cmd, recruitment.DatabaseID, sentTag, timeNow, timeNow); err != nil {
				logger.NewLogger().Error(err.Error())
				return nil, err
			}
		}
	}

	result := model.CreateRecruitmentSuccess{
		RecruitmentEdge: &model.RecruitmentEdge{
			Cursor: utils.GenerateUniqueID("Recruitment", recruitment.DatabaseID),
			Node:   &recruitment,
		},
	}
	return result, nil
}

//* ログインユーザーが作成した募集を取得
func GetViewerRecruitments(ctx context.Context, dbPool *pgxpool.Pool, params search.SearchParams) (*model.RecruitmentConnection, error) {
	cmd := `
		SELECT id, title, type, status, closing_at, created_at, 
		       published_at, prefecture_id, competition_id
		FROM recruitments 
		WHERE user_id = $1
		AND ($2 OR id < $3)
		ORDER BY id DESC
		LIMIT $4
	`
	viewer := user.GetViewer(ctx)
	rows, err := dbPool.Query(
		ctx, cmd,
		viewer.DatabaseID, !params.UseAfter, params.After, params.NumRows,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	connection := model.RecruitmentConnection{
		PageInfo: &model.PageInfo{},
	}
	for rows.Next() {
		var recruitment model.Recruitment
		err := rows.Scan(
			&recruitment.DatabaseID, &recruitment.Title, &recruitment.Type, &recruitment.Status, &recruitment.ClosingAt, &recruitment.CreatedAt,
			&recruitment.PublishedAt, &recruitment.PrefectureID, &recruitment.CompetitionID,
		)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		recruitment.Type = model.Type(strings.ToUpper(recruitment.Type.String()))
		recruitment.Status = model.Status(strings.ToUpper(recruitment.Status.String()))
		connection.Edges = append(connection.Edges, &model.RecruitmentEdge{
			Cursor: utils.GenerateUniqueID("Recruitment", recruitment.DatabaseID),
			Node:   &recruitment,
		})
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	if len(connection.Edges) > 0 {
		lastEdge := connection.Edges[len(connection.Edges)-1]

		cmd = `
		  SELECT COUNT(DISTINCT r.id)
			FROM (
				SELECT *
				FROM recruitments
				WHERE id < $1
				AND user_id = $2
				ORDER BY id DESC
			) as r 
			LIMIT 1
		`
		row := dbPool.QueryRow(
			ctx, cmd,
			lastEdge.Node.DatabaseID, viewer.DatabaseID,
		)

		var count int
		err := row.Scan(&count)
		if err != nil {
			logger.NewLogger().Error(err.Error())
			return nil, err
		}

		var isNextPage bool
		if count > 0 {
			isNextPage = true
		}

		var pageInfo model.PageInfo
		pageInfo.EndCursor = &lastEdge.Cursor
		pageInfo.HasNextPage = isNextPage

		connection.PageInfo = &pageInfo
	}
	return &connection, nil
}

//* IDで募集を取得
func GetRecruitment(ctx context.Context, dbPool *pgxpool.Pool, id int) (*model.Recruitment, error) {
	cmd := `
	  SELECT id, title, type, status, detail, start_at, closing_at, venue, location_lat, 
		       location_lng, user_id, prefecture_id, competition_id, published_at, created_at
		FROM recruitments
		WHERE id = $1
	`

	row := dbPool.QueryRow(ctx, cmd, id)

	var recruitment model.Recruitment
	err := row.Scan(&recruitment.DatabaseID, &recruitment.Title, &recruitment.Type, &recruitment.Status,
		&recruitment.Detail, &recruitment.StartAt, &recruitment.ClosingAt, &recruitment.Venue, &recruitment.LocationLat, &recruitment.LocationLng,
		&recruitment.UserID, &recruitment.PrefectureID, &recruitment.CompetitionID, &recruitment.PublishedAt, &recruitment.CreatedAt,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	recruitment.Status = model.Status(strings.ToUpper(string(recruitment.Status)))
	recruitment.Type = model.Type(strings.ToUpper(string(recruitment.Type)))

	return &recruitment, nil
}

//* 応募済みの募集を取得
func GetAppliedRecruitments(ctx context.Context, dbPool *pgxpool.Pool) ([]*model.Recruitment, error) {
	viewer := user.GetViewer(ctx)

	cmd := `
	  SELECT r.id, r.title, r.type
		FROM recruitments AS r
		INNER JOIN applicants AS a
		ON r.id = a.recruitment_id
		WHERE a.user_id = $1
		ORDER BY a.created_at DESC
	`

	rows, err := dbPool.Query(ctx, cmd, viewer.ID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var recruitments []*model.Recruitment
	for rows.Next() {
		var recruitment model.Recruitment
		err := rows.Scan(&recruitment.ID, &recruitment.Title, &recruitment.Type)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}

		recruitment.Type = model.Type(strings.ToUpper(string(recruitment.Type)))
		recruitments = append(recruitments, &recruitment)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return recruitments, nil
}

//* 全ての募集を取得
func GetRecruitments(ctx context.Context, dbPool *pgxpool.Pool, params search.SearchParams) (*model.RecruitmentConnection, error) {
	cmd := `
		SELECT r.id, r.title, r.type, r.status, r.updated_at, r.closing_at, r.start_at, r.published_at, r.prefecture_id, r.user_id, r.competition_id
		FROM 
			(
				SELECT id, title, type, status, updated_at, closing_at, prefecture_id, user_id, competition_id, published_at, start_at
				FROM recruitments 
				WHERE status = 'published'
				AND ($1 OR id < $2)
				ORDER BY id DESC
				LIMIT $3
			) AS r
		ORDER BY r.id DESC
	`
	rows, err := dbPool.Query(
		ctx, cmd,
		!params.UseAfter, params.After, params.NumRows,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	recConnection := model.RecruitmentConnection{
		PageInfo: &model.PageInfo{},
	}
	for rows.Next() {
		var recruitment model.Recruitment
		if err := rows.Scan(&recruitment.DatabaseID, &recruitment.Title, &recruitment.Type, &recruitment.Status,
			&recruitment.UpdatedAt, &recruitment.ClosingAt, &recruitment.StartAt, &recruitment.PublishedAt,
			&recruitment.PrefectureID, &recruitment.UserID, &recruitment.CompetitionID,
		); err != nil {
			logger.NewLogger().Error(err.Error())
		}
		recConnection.Edges = append(recConnection.Edges, &model.RecruitmentEdge{
			Cursor: utils.GenerateUniqueID("Recruitment", recruitment.DatabaseID),
			Node:   &recruitment,
		})
	}

	if err := rows.Err(); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	if len(recConnection.Edges) > 0 {
		cmd = `
			SELECT COUNT(DISTINCT r.id)
			FROM 
				(
					SELECT id FROM recruitments
					WHERE status = 'published'
					AND id < $1
					ORDER BY id DESC
				) AS r
			LIMIT 1
		`
		lastEdge := recConnection.Edges[len(recConnection.Edges)-1]
		row := dbPool.QueryRow(ctx, cmd, lastEdge.Node.DatabaseID)
		var count int
		err = row.Scan(&count)
		if err != nil {
			logger.NewLogger().Error(err.Error())
			return nil, err
		}
		if count > 0 {
			recConnection.PageInfo.HasNextPage = true
		}
		recConnection.PageInfo.EndCursor = &lastEdge.Cursor
	}
	return &recConnection, nil
}

//* ストックしている募集を取得
func GetStockedRecruitments(ctx context.Context, dbPool *pgxpool.Pool, params search.SearchParams) (*model.RecruitmentConnection, error) {
	viewer := user.GetViewer(ctx)

	cmd := `
		SELECT r.id, r.title, r.closing_at, r.user_id
		FROM recruitments AS r
		INNER JOIN stocks AS s 
			ON r.id = s.recruitment_id
		WHERE r.status = 'published'
		AND s.user_id = $1
		AND ($2 OR s.id < (
												SELECT s.id
												FROM stocks AS s 
												WHERE s.recruitment_id = $3
												AND s.user_id = $4
											))
		ORDER BY s.id DESC
		LIMIT $5
	`

	rows, err := dbPool.Query(
		ctx, cmd,
		viewer.DatabaseID, !params.UseAfter, params.After, viewer.DatabaseID, params.NumRows,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	connection := model.RecruitmentConnection{
		PageInfo: &model.PageInfo{},
		Edges:    []*model.RecruitmentEdge{},
	}

	for rows.Next() {
		var recruitment model.Recruitment
		err := rows.Scan(&recruitment.DatabaseID, &recruitment.Title, &recruitment.ClosingAt, &recruitment.UserID)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		connection.Edges = append(connection.Edges, &model.RecruitmentEdge{
			Cursor: utils.GenerateUniqueID("Recruitment", recruitment.DatabaseID),
			Node:   &recruitment,
		})
	}

	if len(connection.Edges) > 0 {
		lastEdge := connection.Edges[len(connection.Edges)-1]
		cmd = `
			SELECT COUNT(DISTINCT r.id)
			FROM 
			(
				SELECT r.id
				FROM recruitments AS r
				INNER JOIN stocks AS s 
					ON r.id = s.recruitment_id
				WHERE s.user_id = $1
				AND s.id < (
											SELECT s.id
											FROM stocks AS s 
											WHERE s.recruitment_id = $2
											AND s.user_id = $3
										)
				ORDER BY s.id DESC
				LIMIT 1
			) as r
		`
		row := dbPool.QueryRow(
			ctx, cmd,
			viewer.DatabaseID, lastEdge.Node.DatabaseID, viewer.DatabaseID,
		)

		var count int
		err = row.Scan(&count)
		if err != nil {
			logger.NewLogger().Error(err.Error())
			return nil, err
		}
		if count > 0 {
			connection.PageInfo.HasNextPage = true
		}
		connection.PageInfo.EndCursor = &lastEdge.Cursor
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &connection, nil
}

//* ユーザーが作成した募集を取得
func GetUserRecruitments(ctx context.Context, dbPool *pgxpool.Pool, userID int, params search.SearchParams) (*model.RecruitmentConnection, error) {
	cmd := `
	  SELECT id, title, type, prefecture_id, competition_id, closing_at, start_at
		FROM recruitments
		WHERE status = 'published'
		AND user_id = $1
		AND ($2 OR id < $3)
		ORDER BY id DESC
		LIMIT $4
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

	connection := model.RecruitmentConnection{
		PageInfo: &model.PageInfo{},
	}
	for rows.Next() {
		var recruitment model.Recruitment
		err := rows.Scan(
			&recruitment.DatabaseID, &recruitment.Title, &recruitment.Type, &recruitment.PrefectureID,
			&recruitment.CompetitionID, &recruitment.ClosingAt, &recruitment.StartAt,
		)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		connection.Edges = append(connection.Edges, &model.RecruitmentEdge{
			Cursor: utils.GenerateUniqueID("Recruitment", recruitment.DatabaseID),
			Node:   &recruitment,
		})
	}

	if len(connection.Edges) > 0 {
		lastEdge := connection.Edges[len(connection.Edges)-1]
		cmd = `
			SELECT COUNT(DISTINCT r.id)
			FROM (
				SELECT id
				FROM recruitments
				WHERE status = 'published'
				AND user_id = $1
				AND id < $2
				ORDER BY id DESC
				LIMIT 1
			) as r
		`
		row := dbPool.QueryRow(
			ctx, cmd,
			userID, lastEdge.Node.DatabaseID,
		)

		var count int
		err := row.Scan(&count)
		if err != nil {
			logger.NewLogger().Error(err.Error())
			return nil, err
		}
		if count > 0 {
			connection.PageInfo.HasNextPage = true
		}
		connection.PageInfo.EndCursor = &lastEdge.Cursor
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &connection, nil
}

//* 募集の更新
func (i *RecruitmentInput) UpdateRecruitment(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID int) (model.UpdateRecruitmentResult, error) {
	cmd := `
	  UPDATE recruitments
		SET (title, competition_id, type, detail, prefecture_id, venue, 
			  closing_at, start_at, location_lat, location_lng, updated_at, status, published_at) = 
				($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, CASE
																														  WHEN published_at IS NULL 
																															THEN $13
																															ELSE published_at
																														END
				)
		WHERE id = $14
		AND user_id = $15
		RETURNING id, title, status, created_at, published_at, competition_id, 
		          user_id, closing_at, type, prefecture_id, location_lat, location_lng,
							venue, detail, start_at
	`
	viewer := user.GetViewer(ctx)
	timeNow := time.Now().Local()

	var publishedAt *time.Time
	if i.Status == model.StatusPublished {
		publishedAt = &timeNow
	} else {
		publishedAt = nil
	}

	row := dbPool.QueryRow(
		ctx, cmd,
		i.Title, i.CompetitionID, i.Type, i.Detail, i.PrefectureID, i.Venue,
		i.ClosingAt, i.StartAt, i.LocationLat, i.LocationLng, time.Now().Local(), strings.ToLower(string(i.Status)),
		publishedAt, recruitmentID, viewer.DatabaseID,
	)

	var recruitment model.Recruitment
	if err := row.Scan(
		&recruitment.DatabaseID, &recruitment.Title, &recruitment.Status, &recruitment.CreatedAt,
		&recruitment.PublishedAt, &recruitment.CompetitionID, &recruitment.UserID, &recruitment.ClosingAt,
		&recruitment.Type, &recruitment.PrefectureID, &recruitment.LocationLat, &recruitment.LocationLng, &recruitment.Venue, &recruitment.Detail, &recruitment.StartAt,
	); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	//* 募集に紐づいているタグを取得する
	recruitmentTags, err := tag.GetTagsByRecruitmentID(ctx, dbPool, recruitment.DatabaseID)
	if err != nil {
		return nil, err
	}

	//* チェックを外されたタグ(削除するもの)
	var removeTags []int
	//* 削除するタグを見つける処理
	for _, recruitmentTag := range recruitmentTags {
		found := false
		//* 送られてきたタグの中に今のタグがあるか
		for _, sentTagID := range i.TagIDs {
			if recruitmentTag.DatabaseID == sentTagID {
				found = true
			}
		}
		if !found {
			removeTags = append(removeTags, recruitmentTag.DatabaseID)
		}
	}

	//* 新しく付与するタグを見つける処理
	var addTags []int
	for _, sentTag := range i.TagIDs {
		found := false
		for _, recruitmentTag := range recruitmentTags {
			if sentTag == recruitmentTag.DatabaseID {
				found = true
			}
		}
		if !found {
			addTags = append(addTags, sentTag)
		}
	}

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	//* タグを追加
	for _, addTag := range addTags {
		if err := tag.AddRecruitmentTag(ctx, tx, addTag, recruitment.DatabaseID); err != nil {
			if err := tx.Rollback(ctx); err != nil {
				logger.NewLogger().Error(err.Error())
			}
			return nil, err
		}
	}

	//* タグを削除
	for _, removeTag := range removeTags {
		if err := tag.RemoveRecruitmentTag(ctx, tx, removeTag, recruitmentID); err != nil {
			if err := tx.Rollback(ctx); err != nil {
				logger.NewLogger().Error(err.Error())
			}
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	result := model.UpdateRecruitmentSuccess{
		Recruitment: &recruitment,
	}
	return result, nil
}

func DeleteRecruitment(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID int) (*model.DeleteRecruitmentResult, error) {
	cmd := `
	  DELETE FROM recruitments 
		WHERE id = $1 
		AND user_id = $2 
		RETURNING id
	`
	viewer := user.GetViewer(ctx)
	row := dbPool.QueryRow(ctx, cmd, recruitmentID, viewer.DatabaseID)

	var ID int
	if err := row.Scan(&ID); err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	success := model.DeleteRecruitmentResult{
		DeletedRecruitmentID: utils.GenerateUniqueID("Recruitment", ID),
	}
	return &success, nil
}
