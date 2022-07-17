package recruitment

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/search"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/rs/xid"
)

type Recruitment struct {
	Title         string
	Type          model.Type
	Venue         *string
	StartAt       *time.Time
	Detail        *string
	LocationLat   *float64
	LocationLng   *float64
	Status        model.Status
	ClosingAt     *time.Time
	CompetitionID string
	PrefectureID  *string
	Tags          []*model.RecruitmentTagInput
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

func (r Recruitment) RecruitmentValidate() error {
	return validation.ValidateStruct(&r,
		validation.Field(
			&r.Title,
			validation.Required.Error("タイトルを入力してください"),
			validation.RuneLength(1, 60).Error("タイトルは1文字以上60文字以内で入力してください"),
		),
		validation.Field(
			&r.CompetitionID,
			validation.When(r.Status == model.StatusPublished,
				validation.Required.Error("募集競技を選択してください"),
			),
		),
		validation.Field(
			&r.Type,
			validation.In(
				model.TypeOpponent,
				model.TypePersonal,
				model.TypeMember,
				model.TypeJoin,
				model.TypeOther,
			),
		),
		validation.Field(
			&r.Detail,
			validation.When(r.Status == model.StatusPublished,
				validation.Required.Error("募集の詳細を入力してください"),
				validation.RuneLength(1, 10000).Error("募集の詳細は10000文字以内で入力してください"),
			).Else(validation.RuneLength(0, 10000).Error("募集の詳細は10000文字以内で入力してください")),
		),
		validation.Field(
			&r.PrefectureID,
			validation.When(r.Status == model.StatusPublished,
				validation.Required.Error("募集エリアを選択してください"),
			),
		),
		validation.Field(
			&r.Venue,
			validation.When(r.Status == model.StatusPublished,
				validation.When(r.Type == model.TypeOpponent || r.Type == model.TypePersonal,
					validation.Required.Error("会場名を入力してください"),
				),
			),
		),
		validation.Field(
			&r.StartAt,
			validation.When(r.Status == model.StatusPublished,
				validation.When(r.Type == model.TypeOpponent || r.Type == model.TypePersonal,
					validation.By(beforeNowStart),
					validation.Required.Error("開催日時を設定してください"),
				),
			),
		),
		validation.Field(
			&r.ClosingAt,
			validation.When(r.Status == model.StatusPublished,
				validation.Required.Error("募集期限を設定してください"),
				validation.When(r.Type == model.TypeOpponent || r.Type == model.TypePersonal,
					validation.By(beforeNowClosing),
					validation.By(checkWithinTheDeadline(*r.StartAt)),
				),
			),
		),
	)
}

func (r *Recruitment) CreateRecruitment(ctx context.Context, dbPool *pgxpool.Pool) (*model.RecruitmentPayload, error) {
	timeNow := time.Now().Local()

	var publishedAt *time.Time
	if r.Status == model.StatusPublished {
		publishedAt = &timeNow
	} else {
		publishedAt = nil
	}

	cmd := `
	  INSERT INTO recruitments 
		  (title, competition_id, type, detail, prefecture_id, venue, start_at, closing_at, location_lat, location_lng, status, user_id, created_at, updated_at, published_at)
		VALUES
		  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, title, status, created_at, published_at, competition_id, user_id, closing_at, type, prefecture_id
		`

	viewer := auth.ForContext(ctx)
	row := dbPool.QueryRow(
		ctx, cmd,
		r.Title, r.CompetitionID, strings.ToLower(string(r.Type)), r.Detail, r.PrefectureID, r.Venue, r.StartAt,
		r.ClosingAt, r.LocationLat, r.LocationLng, strings.ToLower(string(r.Status)), viewer.DatabaseID, timeNow, timeNow, publishedAt,
	)

	var payload model.RecruitmentPayload
	var recruitment model.Recruitment
	err := row.Scan(&recruitment.DatabaseID, &recruitment.Title, &recruitment.Status, &recruitment.CreatedAt, &recruitment.PublishedAt,
		&recruitment.CompetitionID, &recruitment.UserID, &recruitment.ClosingAt, &recruitment.Type, &recruitment.PrefectureID,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	var recruitmentEdge model.RecruitmentEdge

	recruitment.Status = model.Status(strings.ToUpper(recruitment.Status.String()))
	recruitment.Type = model.Type(strings.ToUpper(recruitment.Type.String()))
	recruitmentEdge.Node = &recruitment
	recruitmentEdge.Cursor = utils.GenerateUniqueID("Recruitment", recruitment.DatabaseID)
	payload.FeedbackRecruitmentEdge = &recruitmentEdge

	cmd = "INSERT INTO tags (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, name"

	if len(r.Tags) > 0 {

		var newTags []*model.Tag
		for _, tag := range r.Tags {
			if tag.IsNew {
				row := dbPool.QueryRow(ctx, cmd, xid.New().String(), tag.Name, timeNow, timeNow)

				var tag model.Tag
				err := row.Scan(&tag.ID, &tag.Name)
				if err != nil {
					logger.NewLogger().Error(err.Error())
				}
				newTags = append(newTags, &tag)
			} else {
				newTags = append(newTags, &model.Tag{ID: tag.ID, Name: tag.Name})
			}
		}

		cmd = "INSERT INTO recruitment_tags (id, recruitment_id, tag_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"

		for _, tag := range newTags {
			if _, err := dbPool.Exec(ctx, cmd, xid.New().String(), recruitment.ID, tag.ID, timeNow, timeNow); err != nil {
				logger.NewLogger().Error(err.Error())
			}
		}
	}

	fmt.Println(*payload.FeedbackRecruitmentEdge.Node)
	return &payload, nil
}

func GetViewerRecruitments(ctx context.Context, dbPool *pgxpool.Pool, params search.SearchParams) (*model.RecruitmentConnection, error) {
	viewer := auth.ForContext(ctx)

	cmd := `
		SELECT r.id, r.title, r.type, r.status, r.closing_at, r.created_at, 
		       r.published_at, r.prefecture_id, r.competition_id
		FROM 
			(
				SELECT *
				FROM recruitments 
				WHERE ($1 OR id < $2)
				ORDER BY id DESC
				LIMIT $3
			) AS r
		WHERE r.user_id = $4
		ORDER BY r.id DESC
	`

	rows, err := dbPool.Query(
		ctx, cmd,
		!params.UseAfter, params.After, params.NumRows, viewer.DatabaseID,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var connection model.RecruitmentConnection
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
		endCursor := connection.Edges[len(connection.Edges)-1].Cursor

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
			utils.DecodeUniqueID(endCursor), viewer.DatabaseID,
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
		pageInfo.EndCursor = &endCursor
		pageInfo.HasNextPage = isNextPage

		connection.PageInfo = &pageInfo
	} else {
		connection.PageInfo = &model.PageInfo{}
	}

	return &connection, nil
}

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

func GetAppliedRecruitments(ctx context.Context, dbPool *pgxpool.Pool) ([]*model.Recruitment, error) {
	viewer := auth.ForContext(ctx)

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

func GetRecruitments(ctx context.Context, dbPool *pgxpool.Pool, params search.SearchParams) (*model.RecruitmentConnection, error) {
	var sort string
	if params.UseBefore {
		sort = "ASC"
	} else {
		sort = "DESC"
	}

	cmd := fmt.Sprintf(`
		SELECT r.id, r.title, r.type, r.status, r.updated_at, r.closing_at, r.published_at, r.prefecture_id, r.user_id, r.competition_id
		FROM 
			(
				SELECT id, title, type, status, updated_at, closing_at, prefecture_id, user_id, competition_id, published_at
				FROM recruitments 
				WHERE status = 'published'
				AND ($1 OR id < $2)
				ORDER BY id %s
				LIMIT $3
			) AS r
		ORDER BY r.id DESC
	`, sort)

	rows, err := dbPool.Query(
		ctx, cmd,
		!params.UseAfter, params.After, params.NumRows,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var recConnection model.RecruitmentConnection

	for rows.Next() {
		var recruitment model.Recruitment

		err := rows.Scan(&recruitment.DatabaseID, &recruitment.Title, &recruitment.Type, &recruitment.Status,
			&recruitment.UpdatedAt, &recruitment.ClosingAt, &recruitment.PublishedAt, &recruitment.PrefectureID, &recruitment.UserID, &recruitment.CompetitionID)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}

		recruitment.Type = model.Type(strings.ToUpper(string(recruitment.Type)))
		recruitment.Status = model.Status(strings.ToUpper(string(recruitment.Status)))
		recConnection.Edges = append(recConnection.Edges, &model.RecruitmentEdge{
			Cursor: utils.GenerateUniqueID("Recruitment", recruitment.DatabaseID),
			Node:   &recruitment,
		})
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return &model.RecruitmentConnection{}, err
	}

	if len(recConnection.Edges) > 0 {
		endCursor := recConnection.Edges[len(recConnection.Edges)-1].Cursor

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

		isNext, err := search.NextPageExists(ctx, dbPool, endCursor, cmd)
		if err != nil {
			return &model.RecruitmentConnection{}, err
		}

		var pageInfo model.PageInfo

		pageInfo.HasNextPage = isNext
		pageInfo.EndCursor = &endCursor

		recConnection.PageInfo = &pageInfo
	} else {
		recConnection.PageInfo = &model.PageInfo{}
	}

	return &recConnection, nil
}

func GetStockedRecruitments(ctx context.Context, dbPool *pgxpool.Pool) ([]*model.Recruitment, error) {
	viewer := auth.ForContext(ctx)

	cmd := `
		SELECT r.id, r.title, r.closing_at, r.user_id
		FROM recruitments AS r
		INNER JOIN stocks AS s
		ON s.recruitment_id = r.id
		WHERE s.user_id = $1
		ORDER BY s.created_at DESC
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
		err := rows.Scan(&recruitment.ID, &recruitment.Title, &recruitment.ClosingAt)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		recruitments = append(recruitments, &recruitment)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return recruitments, nil
}

func (r *Recruitment) UpdateRecruitment(ctx context.Context, dbPool *pgxpool.Pool, recID string) (*model.Recruitment, error) {
	viewer := auth.ForContext(ctx)
	if viewer == nil {
		return nil, errors.New("ログインしてください")
	}

	timeNow := time.Now().Local()

	var publishedAt *time.Time

	if r.Status == model.StatusPublished {
		publishedAt = &timeNow
	} else {
		publishedAt = nil
	}

	cmd := `
	  UPDATE recruitments AS r
		SET title = $1, competition_id = $2, type = $3, detail = $4, prefecture_id = $5, venue = $6,
		    closing_at = $7, start_at = $8, location_lat = $9, location_lng = $10, updated_at = $11, status = $12,
				published_at = CASE 
												 WHEN r.published_at IS NULL 
												 THEN $13
												 ELSE r.published_at
											 END 
		WHERE r.id = $14
		AND r.user_id = $15
		RETURNING id
	`

	row := dbPool.QueryRow(
		ctx, cmd,
		r.Title, r.CompetitionID, r.Type, r.Detail, r.PrefectureID, r.Venue,
		r.ClosingAt, r.StartAt, r.LocationLat, r.LocationLng, time.Now().Local(), strings.ToLower(string(r.Status)), publishedAt, recID, viewer.ID,
	)

	var ID string
	err := row.Scan(&ID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	cmd = `
	  SELECT t.id, t.name 
		FROM tags AS t
		INNER JOIN recruitment_tags AS r_t
		  ON r_t.tag_id = t.id
		WHERE r_t.recruitment_id = $1
	`

	rows, err := dbPool.Query(ctx, cmd, ID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var currentTags []*model.Tag // 更新する募集に既についているタグ
	for rows.Next() {
		var tag model.Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		currentTags = append(currentTags, &tag)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	cmd = "INSERT INTO tags (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, name"

	var sentTags []*model.Tag // フォームから送られてきたタグ
	for _, tag := range r.Tags {
		if tag.IsNew {
			row := dbPool.QueryRow(ctx, cmd, xid.New().String(), tag.Name, timeNow, timeNow)

			var tag model.Tag
			err := row.Scan(&tag.ID, &tag.Name)
			if err != nil {
				logger.NewLogger().Error(err.Error())
			}

			sentTags = append(sentTags, &tag)
		} else {
			sentTags = append(sentTags, &model.Tag{ID: tag.ID, Name: tag.Name})
		}
	}

	var oldTags []*model.Tag // チェックを外されたタグ(削除するもの)

	// 削除するタグを見つける処理
	for _, currentTag := range currentTags {
		found := false
		for _, sentTag := range sentTags {
			if currentTag.Name == sentTag.Name {
				found = true
			}
		}
		if !found {
			oldTags = append(oldTags, currentTag)
		}
	}

	tx, err := dbPool.Begin(ctx)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	cmd = `
	  INSERT INTO recruitment_tags 
		  (id, recruitment_id, tag_id, created_at, updated_at) 
		VALUES 
		  ($1, $2, $3, $4, $5)
		ON CONFLICT 
		  ON CONSTRAINT 
			  recruitment_tags_recruitment_id_tag_id_key
		DO UPDATE SET updated_at = $6
	`

	for _, tag := range sentTags {
		timeNow := time.Now().Local()
		if _, err := tx.Exec(ctx, cmd, xid.New().String(), ID, tag.ID, timeNow, timeNow, timeNow); err != nil {
			logger.NewLogger().Error(err.Error())
			err = tx.Rollback(ctx)
			if err != nil {
				logger.NewLogger().Error(err.Error())
			}
			return nil, err
		}
	}

	cmd = `
	  DELETE FROM recruitment_tags AS r_t
		WHERE r_t.recruitment_id = $1 AND r_t.tag_id = $2
	`

	for _, tag := range oldTags {
		if _, err := tx.Exec(ctx, cmd, ID, tag.ID); err != nil {
			logger.NewLogger().Error(err.Error())
			err = tx.Rollback(ctx)
			if err != nil {
				logger.NewLogger().Error(err.Error())
			}
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		err = tx.Rollback(ctx)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		return nil, err
	}

	cmd = `
	  SELECT r.id, r.title, r.status
		FROM recruitments AS r 
		WHERE r.id = $1
	`

	row = dbPool.QueryRow(ctx, cmd, ID)

	var recruitment model.Recruitment
	err = row.Scan(&recruitment.ID, &recruitment.Title, &recruitment.Status)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &recruitment, nil
}

func DeleteRecruitment(ctx context.Context, dbPool *pgxpool.Pool, recruitmentID string) (*model.DeleteRecruitmentPayload, error) {
	viewer := auth.ForContext(ctx)

	cmd := "DELETE FROM recruitments AS r WHERE r.id = $1 AND r.user_id = $2 RETURNING id"
	row := dbPool.QueryRow(ctx, cmd, utils.DecodeUniqueID(recruitmentID), viewer.DatabaseID)

	var recruitment model.Recruitment
	err := row.Scan(&recruitment.DatabaseID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &model.DeleteRecruitmentPayload{DeletedRecruitmentID: utils.GenerateUniqueID("Recruitment", recruitment.DatabaseID)}, nil
}
