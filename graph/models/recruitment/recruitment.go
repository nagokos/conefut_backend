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
	"github.com/nagokos/connefut_backend/graph/models/prefecture"
	"github.com/nagokos/connefut_backend/graph/models/search"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/rs/xid"
)

type Recruitment struct {
	Title         string
	Type          model.Type
	Place         *string
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
				model.TypeIndividual,
				model.TypeMember,
				model.TypeJoining,
				model.TypeOthers,
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
			&r.Place,
			validation.When(r.Status == model.StatusPublished,
				validation.When(r.Type == model.TypeOpponent || r.Type == model.TypeIndividual,
					validation.Required.Error("会場名を入力してください"),
				),
			),
		),
		validation.Field(
			&r.StartAt,
			validation.When(r.Status == model.StatusPublished,
				validation.When(r.Type == model.TypeOpponent || r.Type == model.TypeIndividual,
					validation.By(beforeNowStart),
					validation.Required.Error("開催日時を設定してください"),
				),
			),
		),
		validation.Field(
			&r.ClosingAt,
			validation.When(r.Status == model.StatusPublished,
				validation.Required.Error("募集期限を設定してください"),
				validation.When(r.Type == model.TypeOpponent || r.Type == model.TypeIndividual,
					validation.By(beforeNowClosing),
					validation.By(checkWithinTheDeadline(*r.StartAt)),
				),
			),
		),
	)
}

func (r *Recruitment) CreateRecruitment(ctx context.Context, dbPool *pgxpool.Pool) (*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return &model.Recruitment{}, errors.New("ログインしてください")
	}

	timeNow := time.Now().Local()

	var published_at *time.Time
	if r.Status == model.StatusPublished {
		published_at = &timeNow
	} else {
		published_at = nil
	}

	cmd := `
	  INSERT INTO recruitments 
		  (id, title, competition_id, type, detail, prefecture_id, place, start_at, closing_at, location_lat, location_lng, status, user_id, created_at, updated_at, published_at)
		VALUES
		  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, title, status 
		`

	row := dbPool.QueryRow(
		ctx, cmd,
		xid.New().String(), r.Title, r.CompetitionID, r.Type, r.Detail, r.PrefectureID, r.Place, r.StartAt,
		r.ClosingAt, r.LocationLat, r.LocationLng, strings.ToLower(string(r.Status)), currentUser.ID, timeNow, timeNow, published_at,
	)

	var recruitment model.Recruitment
	err := row.Scan(&recruitment.ID, &recruitment.Title, &recruitment.Status)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	cmd = "INSERT INTO tags (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, name"

	if len(r.Tags) != 0 {
		tx, err := dbPool.Begin(ctx)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		defer tx.Rollback(ctx)

		var newTags []*model.Tag
		for _, tag := range r.Tags {
			if tag.IsNew {
				row := tx.QueryRow(ctx, cmd, xid.New().String(), tag.Name, timeNow, timeNow)

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
			if _, err := tx.Exec(ctx, cmd, xid.New().String(), recruitment.ID, tag.ID, timeNow, timeNow); err != nil {
				logger.NewLogger().Error(err.Error())
			}
		}

		if err := tx.Commit(ctx); err != nil {
			logger.NewLogger().Error(err.Error())
		}
	}

	return &recruitment, nil
}

func GetCurrentUserRecruitments(ctx context.Context, dbPool *pgxpool.Pool) ([]*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)

	cmd := `
	  SELECT r.id, r.title, r.type, r.status, r.closing_at, r.created_at, r.published_at
		FROM recruitments AS r
		WHERE r.user_id = $1
		ORDER BY r.status DESC, r.updated_at DESC
		`

	rows, err := dbPool.Query(ctx, cmd, currentUser.ID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var recruitments []*model.Recruitment
	for rows.Next() {
		var recruitment model.Recruitment
		err := rows.Scan(&recruitment.ID, &recruitment.Title, &recruitment.Type, &recruitment.Status, &recruitment.ClosingAt,
			&recruitment.CreatedAt, &recruitment.PublishedAt,
		)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		recruitment.Status = model.Status(strings.ToUpper(string(recruitment.Status)))
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

func GetRecruitment(ctx context.Context, dbPool *pgxpool.Pool, recID string) (*model.Recruitment, error) {
	cmd := `
		SELECT r.id, r.title, r.type, r.status, r.detail, r.start_at, r.closing_at, r.place, r.location_lat, r.location_lng,
		       c.id, c.name,
					 p.id, p.name,
					 u.id, u.name, u.avatar
		FROM recruitments AS r
		LEFT OUTER JOIN prefectures AS p 
		ON r.prefecture_id = p.id
		INNER JOIN competitions AS c
			ON r.competition_id = c.id
		INNER JOIN users AS u 
			ON r.user_id = u.id
		WHERE r.id = $1
	`

	row := dbPool.QueryRow(ctx, cmd, recID)

	var recruitment model.Recruitment
	var competition model.Competition
	var nullablePrefecture prefecture.NullablePrefecture
	var user model.User
	err := row.Scan(&recruitment.ID, &recruitment.Title, &recruitment.Type, &recruitment.Status,
		&recruitment.Detail, &recruitment.StartAt, &recruitment.ClosingAt, &recruitment.Place, &recruitment.LocationLat, &recruitment.LocationLng,
		&competition.ID, &competition.Name, &nullablePrefecture.ID, &nullablePrefecture.Name, &user.ID, &user.Name, &user.Avatar,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	recruitment.Competition = &model.Competition{ID: competition.ID, Name: competition.Name}
	recruitment.User = &model.User{ID: user.ID, Name: user.Name, Avatar: user.Avatar}
	recruitment.Status = model.Status(strings.ToUpper(string(recruitment.Status)))
	recruitment.Type = model.Type(strings.ToUpper(string(recruitment.Type)))

	var defaultPrefecture prefecture.NullablePrefecture
	if nullablePrefecture != defaultPrefecture {
		recruitment.Prefecture = &model.Prefecture{ID: *nullablePrefecture.ID, Name: *nullablePrefecture.Name}
	}

	cmd = `
		SELECT t.id, t.name
		FROM tags AS t
		INNER JOIN recruitment_tags AS r_t
			ON r_t.tag_id = t.id
		WHERE r_t.recruitment_id = $1
	`

	rows, err := dbPool.Query(ctx, cmd, recruitment.ID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag model.Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		recruitment.Tags = append(recruitment.Tags, &tag)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &recruitment, nil
}

func GetAppliedRecruitments(ctx context.Context, dbPool *pgxpool.Pool) ([]*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)

	cmd := `
	  SELECT r.id, r.title, r.type, a.created_at AS app_created_at,  u.name AS usr_name, u.avatar AS usr_avatar
		FROM recruitments AS r
		INNER JOIN applicants AS a 
		  ON r.id = a.recruitment_id
		INNER JOIN users AS u 
		  ON u.id = r.user_id
		WHERE a.user_id = $1
		ORDER BY a.created_at DESC
	`

	rows, err := dbPool.Query(ctx, cmd, currentUser.ID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var recruitments []*model.Recruitment
	for rows.Next() {
		var recruitment model.Recruitment
		var applicant model.Applicant
		var user model.User
		err := rows.Scan(&recruitment.ID, &recruitment.Title, &recruitment.Type, &applicant.CreatedAt,
			&user.Name, &user.Avatar,
		)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		recruitment.Applicant = &applicant
		recruitment.User = &user
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
		SELECT r.id, r.title, r.type, r.status, r.updated_at, r.closing_at, r.published_at,
					 u.name, u.avatar,
					 p.name,
					 c.name
		FROM 
			(
				SELECT id, title, type, status, updated_at, closing_at, prefecture_id, user_id, competition_id, published_at
				FROM recruitments 
				WHERE status = $1
				AND ($2 OR competition_id = $3)
				AND ($4 OR id < $5)
				AND ($6 OR id > $7)
				ORDER BY id %s
				LIMIT $8
			) AS r
		INNER JOIN prefectures AS p 
			ON r.prefecture_id = p.id
		INNER JOIN users AS u 
			ON r.user_id = u.id
		INNER JOIN competitions AS c 
			ON r.competition_id = c.id
		ORDER BY r.id DESC
	`, sort)

	rows, err := dbPool.Query(
		ctx, cmd,
		"published", !params.Options.UseCompetition, params.Options.CompetitionID, !params.UseAfter, params.After, !params.UseBefore, params.Before, params.NumRows,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var recConnection model.RecruitmentConnection

	for rows.Next() {
		var recruitment model.Recruitment
		var user model.User
		var prefecture model.Prefecture
		var competition model.Competition

		err := rows.Scan(&recruitment.ID, &recruitment.Title, &recruitment.Type, &recruitment.Status, &recruitment.UpdatedAt, &recruitment.ClosingAt, &recruitment.PublishedAt,
			&user.Name, &user.Avatar, &prefecture.Name, &competition.Name,
		)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}

		recruitment.User = &user
		recruitment.Prefecture = &prefecture
		recruitment.Competition = &competition
		recruitment.Type = model.Type(strings.ToUpper(string(recruitment.Type)))
		recruitment.Status = model.Status(strings.ToUpper(string(recruitment.Status)))
		recConnection.Edges = append(recConnection.Edges, &model.RecruitmentEdge{
			Cursor: recruitment.ID,
			Node:   &recruitment,
		})
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return &model.RecruitmentConnection{}, err
	}

	if len(recConnection.Edges) > 0 {
		startCursor := recConnection.Edges[0].Node.ID
		endCursor := recConnection.Edges[len(recConnection.Edges)-1].Node.ID

		isNext, err := search.NextPageExists(ctx, dbPool, endCursor, params, sort)
		if err != nil {
			return &model.RecruitmentConnection{}, err
		}

		isPrevious, err := search.PreviousPageExists(ctx, dbPool, startCursor, params, sort)
		if err != nil {
			return &model.RecruitmentConnection{}, err
		}

		var pageInfo model.PageInfo

		pageInfo.HasNextPage = isNext
		pageInfo.HasPreviousPage = isPrevious
		pageInfo.StartCursor = startCursor
		pageInfo.EndCursor = endCursor

		recConnection.PageInfo = &pageInfo
	} else {
		recConnection.PageInfo = &model.PageInfo{}
	}

	return &recConnection, nil
}

func GetStockedRecruitments(ctx context.Context, dbPool *pgxpool.Pool) ([]*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)

	cmd := `
		SELECT DISTINCT r.id, r.title, r.closing_at, u.id AS usr_id, u.name AS usr_name, u.avatar AS usr_avatar
		FROM recruitments AS r 
		INNER JOIN stocks AS s 
			ON r.id = s.recruitment_id
		INNER JOIN users AS u 
			ON r.user_id = u.id
		WHERE s.user_id = $1
	`

	rows, err := dbPool.Query(ctx, cmd, currentUser.ID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var recruitments []*model.Recruitment
	for rows.Next() {
		var recruitment model.Recruitment
		var user model.User
		err := rows.Scan(&recruitment.ID, &recruitment.Title, &recruitment.ClosingAt,
			&user.ID, &user.Name, &user.Avatar,
		)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		recruitment.User = &user
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
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return nil, errors.New("ログインしてください")
	}

	timeNow := time.Now().Local()

	var published_at *time.Time

	if r.Status == model.StatusPublished {
		published_at = &timeNow
	} else {
		published_at = nil
	}

	cmd := `
	  UPDATE recruitments AS r
		SET title = $1, competition_id = $2, type = $3, detail = $4, prefecture_id = $5, place = $6,
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
		r.Title, r.CompetitionID, r.Type, r.Detail, r.PrefectureID, r.Place,
		r.ClosingAt, r.StartAt, r.LocationLat, r.LocationLng, time.Now().Local(), strings.ToLower(string(r.Status)), published_at, recID, currentUser.ID,
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
			tx.Rollback(ctx)
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
			tx.Rollback(ctx)
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		tx.Rollback(ctx)
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

func DeleteRecruitment(ctx context.Context, dbPool *pgxpool.Pool, recID string) (*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return &model.Recruitment{}, errors.New("ログインしてください")
	}

	cmd := "DELETE FROM recruitments AS r WHERE r.id = $1 AND r.user_id = $2 RETURNING id, title"
	row := dbPool.QueryRow(ctx, cmd, recID, currentUser.ID)

	var recruitment model.Recruitment
	err := row.Scan(&recruitment.ID, &recruitment.Title)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return &recruitment, err
	}

	return &recruitment, nil
}
