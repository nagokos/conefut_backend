package recruitment

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/auth"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/rs/xid"
)

type Recruitment struct {
	Title         string
	Type          model.Type
	Place         *string
	StartAt       *time.Time
	Content       *string
	LocationLat   *float64
	LocationLng   *float64
	Status        model.Status
	Capacity      *int
	ClosingAt     *time.Time
	CompetitionID *string
	PrefectureID  *string
	Tags          []*model.RecruitmentTagInput
}

func requiredIfUnnecessaryType() validation.RuleFunc {
	return func(v interface{}) error {
		if v == model.TypeUnnecessary {
			return errors.New("募集タイプを選択してください")
		}
		return nil
	}
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
				model.TypeUnnecessary,
				model.TypeOpponent,
				model.TypeIndividual,
				model.TypeMember,
				model.TypeJoining,
				model.TypeOthers,
			),
			validation.When(r.Status == model.StatusPublished,
				validation.By(requiredIfUnnecessaryType()),
			),
		),
		validation.Field(
			&r.Content,
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
			&r.Capacity,
			validation.When(r.Status == model.StatusPublished,
				validation.When(
					r.Type == model.TypeOpponent ||
						r.Type == model.TypeIndividual,
					validation.Required.Error("募集人数は1名以上にしてください"),
					validation.Min(1).Error("募集人数は1名以上にしてください"),
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

func (r *Recruitment) CreateRecruitment(ctx context.Context, dbConnection *sql.DB) (*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return &model.Recruitment{}, errors.New("ログインしてください")
	}

	cmd := `
	  INSERT INTO recruitments 
		  (id, title, competition_id, type, content, prefecture_id, place, capacity, start_at, closing_at, location_lat, location_lng, status, user_id, created_at, updated_at)
		VALUES
		  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, title, status 
		`

	timeNow := time.Now().Local()
	row := dbConnection.QueryRow(
		cmd,
		xid.New().String(), r.Title, r.CompetitionID, r.Type, r.Content, r.PrefectureID, r.Place, r.Capacity, r.StartAt,
		r.ClosingAt, r.LocationLat, r.LocationLng, r.Status, currentUser.ID, timeNow, timeNow,
	)

	var recruitment model.Recruitment
	err := row.Scan(&recruitment.ID, &recruitment.Title, &recruitment.Status)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	if len(r.Tags) != 0 {
		tx, err := dbConnection.Begin()
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		defer tx.Rollback()

		stmt, err := tx.Prepare("INSERT INTO recruitment_tags (id, recruitment_id, tag_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)")
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		defer stmt.Close()

		for _, tag := range r.Tags {
			if _, err := stmt.Exec(xid.New().String(), recruitment.ID, tag.ID, timeNow, timeNow); err != nil {
				logger.NewLogger().Error(err.Error())
			}
		}

		if err := tx.Commit(); err != nil {
			logger.NewLogger().Error(err.Error())
		}
	}

	return &recruitment, nil
}

func GetCurrentUserRecruitments(ctx context.Context, dbConnection *sql.DB) ([]*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)

	cmd := `
	  SELECT r.id, r.title, r.type, r.status
		FROM recruitments AS r
		WHERE r.user_id = $1
		`
	rows, err := dbConnection.Query(cmd, currentUser.ID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var recruitments []*model.Recruitment
	for rows.Next() {
		var recruitment model.Recruitment
		err := rows.Scan(&recruitment.ID, &recruitment.Title, &recruitment.Type, &recruitment.Status)
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

func GetRecruitment(dbConnection *sql.DB, recID string) (*model.Recruitment, error) {
	cmd := `
		SELECT r.id, r.title, r.type, r.status, r.content, r.start_at, r.closing_at, r.place, r.location_lat, r.location_lng,
		       c.id AS comp_id, c.name AS comp_name, 
					 p.id AS pref_id, p.name AS pref_name, 
					 u.id AS usr_id, u.name AS usr_name, u.avatar AS usr_avatar
		FROM recruitments AS r
		LEFT OUTER JOIN competitions AS c
			ON r.competition_id = c.id
		LEFT OUTER JOIN prefectures AS p 
			ON r.prefecture_id = p.id
		INNER JOIN users AS u 
			ON r.user_id = u.id
		WHERE r.id = $1
		ORDER BY id ASC
	`

	row := dbConnection.QueryRow(cmd, recID)

	var recruitment model.Recruitment
	var competition model.Competition
	var prefecture model.Prefecture
	var user model.User
	err := row.Scan(&recruitment.ID, &recruitment.Title, &recruitment.Type, &recruitment.Status,
		&recruitment.Content, &recruitment.StartAt, &recruitment.ClosingAt, &recruitment.Place, &recruitment.LocationLat, &recruitment.LocationLng,
		&competition.ID, &competition.Name, &prefecture.ID, &prefecture.Name, &user.ID, &user.Name, &user.Avatar,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	recruitment.Competition = &competition
	recruitment.Prefecture = &prefecture
	recruitment.User = &user
	recruitment.Status = model.Status(strings.ToUpper(string(recruitment.Status)))
	recruitment.Type = model.Type(strings.ToUpper(string(recruitment.Type)))

	cmd = `
		SELECT t.id, t.name
		FROM tags AS t
		INNER JOIN recruitment_tags AS r_t
			ON r_t.tag_id = t.id
		WHERE r_t.recruitment_id = $1
	`

	rows, err := dbConnection.Query(cmd, recruitment.ID)
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

func GetAppliedRecruitments(ctx context.Context, dbConnection *sql.DB) ([]*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)

	cmd := `
	  SELECT r.id, r.title, r.type, a.created_at AS app_created_at, a.management_status AS app_management_status, u.name AS usr_name, u.avatar AS usr_avatar
		FROM recruitments AS r
		INNER JOIN applicants AS a 
		  ON r.id = a.recruitment_id
		INNER JOIN users AS u 
		  ON u.id = r.user_id
		WHERE r.user_id = $1
	`

	rows, err := dbConnection.Query(cmd, currentUser.ID)
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
			&applicant.ManagementStatus, &user.Name, &user.Avatar,
		)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		recruitment.Applicant = &applicant
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

func GetRecruitments(dbConnection *sql.DB) ([]*model.Recruitment, error) {
	cmd := `
	  SELECT r.id, r.title, r.type, r.updated_at, r.closing_at, u.name AS usr_name, u.avatar AS usr_avatar, p.name AS pref_name
		FROM recruitments AS r
		INNER JOIN prefectures AS p 
		  ON r.prefecture_id = p.id
		INNER JOIN users AS u 
		  ON r.user_id = u.id
	`

	rows, err := dbConnection.Query(cmd)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var recruitments []*model.Recruitment
	for rows.Next() {
		var recruitment model.Recruitment
		var user model.User
		var prefecture model.Prefecture
		err := rows.Scan(&recruitment.ID, &recruitment.Title, &recruitment.Type, &recruitment.UpdatedAt, &recruitment.ClosingAt,
			&user.Name, &user.Avatar, &prefecture.Name,
		)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		recruitment.User = &user
		recruitment.Prefecture = &prefecture
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

func GetStockedRecruitments(ctx context.Context, dbConnection *sql.DB) ([]*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)

	cmd := `
	  SELECT r.id, r.title, r.type, r.status, u.id AS usr_id, u.name AS usr_name, u.avatar AS usr_avatar
		FROM recruitments AS r 
		INNER JOIN stocks AS s 
		  ON s.user_id = $1
		INNER JOIN users AS u 
		  ON r.user_id = u.id
	`

	rows, err := dbConnection.Query(cmd, currentUser.ID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var recruitments []*model.Recruitment
	for rows.Next() {
		var recruitment model.Recruitment
		var user model.User
		err := rows.Scan(&recruitment.ID, &recruitment.Title, &recruitment.Type, &recruitment.Status,
			&user.ID, &user.Name, &user.Avatar,
		)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		recruitment.User = &user
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

func (r *Recruitment) UpdateRecruitment(ctx context.Context, dbConnection *sql.DB, recID string) (*model.Recruitment, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return nil, errors.New("ログインしてください")
	}

	cmd := `
	  UPDATE recruitments AS r
		SET title = $1, competition_id = $2, type = $3, content = $4, prefecture_id = $5, place = $6, capacity = $7, 
		    closing_at = $8, start_at = $9, location_lat = $10, location_lng = $11, updated_at = $12
		WHERE r.id = $13
		AND r.user_id = $14
		RETURNING id
	`

	row := dbConnection.QueryRow(
		cmd,
		r.Title, r.CompetitionID, r.Type, r.Content, r.PrefectureID, r.Place, r.Capacity,
		r.ClosingAt, r.StartAt, r.LocationLat, r.LocationLng, time.Now().Local(), recID, currentUser.ID,
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

	rows, err := dbConnection.Query(cmd, ID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	var currentTags []*model.Tag
	for rows.Next() {
		var tag model.Tag
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		currentTags = append(currentTags, &tag)
	}

	err = row.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	var oldTags []*model.Tag // チェックを外されたタグ(削除するもの)

	// 削除するタグを見つける処理
	for _, currentTag := range currentTags {
		found := false
		for _, sentTag := range r.Tags {
			if currentTag.Name == sentTag.Name {
				found = true
			}
		}
		if !found {
			oldTags = append(oldTags, currentTag)
		}
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

	tx, err := dbConnection.Begin()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	stmt, err := tx.Prepare(cmd)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		tx.Rollback()
		return nil, err
	}

	for _, tag := range r.Tags {
		timeNow := time.Now().Local()
		if _, err := stmt.Exec(xid.New().String(), ID, tag.ID, timeNow, timeNow, timeNow); err != nil {
			logger.NewLogger().Error(err.Error())
			tx.Rollback()
			return nil, err
		}
	}
	stmt.Close()

	cmd = `
	  DELETE FROM recruitment_tags AS r_t
		WHERE r_t.recruitment_id = $1 AND r_t.tag_id = $2
	`

	stmt, err = tx.Prepare(cmd)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		tx.Rollback()
		return nil, err
	}

	for _, tag := range oldTags {
		logger.NewLogger().Info(tag.ID)
		if _, err := stmt.Exec(ID, tag.ID); err != nil {
			logger.NewLogger().Error(err.Error())
			tx.Rollback()
			return nil, err
		}
	}

	stmt.Close()
	tx.Commit()

	cmd = `
	  SELECT r.id, r.title, r.status
		FROM recruitments AS r 
		WHERE r.id = $1
	`

	row = dbConnection.QueryRow(cmd, ID)

	var recruitment model.Recruitment
	err = row.Scan(&recruitment.ID, &recruitment.Title, &recruitment.Status)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &recruitment, nil
}

func DeleteRecruitment(ctx context.Context, dbConnection *sql.DB, recID string) (bool, error) {
	currentUser := auth.ForContext(ctx)
	if currentUser == nil {
		return false, errors.New("ログインしてください")
	}

	cmd := "DELETE FROM recruitments AS r WHERE r.id = $1 AND r.user_id = $2"
	_, err := dbConnection.Exec(cmd, recID, currentUser.ID)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return false, err
	}

	return true, nil
}
