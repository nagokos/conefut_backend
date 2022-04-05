package tag

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
	"github.com/rs/xid"
)

type Tag struct {
	Name string
}

func existsTag() validation.RuleFunc {
	return func(v interface{}) error {

		s := v.(string)
		lower := strings.ToLower(s)
		dbConnection := db.DatabaseConnection()

		cmd := "SELECT COUNT(DISTINCT id) FROM tags WHERE name = $1"
		row := dbConnection.QueryRow(cmd, lower)

		var count int
		err := row.Scan(&count)

		if err != nil {
			logger.NewLogger().Error(err.Error())
			return err
		}

		if count == 1 {
			logger.NewLogger().Error("This tag name is already exists")
			err = errors.New("このタグは既に存在します")
		}

		return err
	}
}

func (t Tag) CreateTagValidate() error {
	return validation.ValidateStruct(&t,
		validation.Field(
			&t.Name,
			validation.Required.Error("タグ名を入力してください"),
			validation.By(existsTag()),
		),
	)
}

func GetTags(dbConnection *sql.DB) ([]*model.Tag, error) {
	var tags []*model.Tag

	cmd := "SELECT id, name FROM tags"
	rows, err := dbConnection.Query(cmd)
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
		tags = append(tags, &tag)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return tags, nil
}

func GetRecruitmentTags(dbConnection *sql.DB, recId string) ([]*model.Tag, error) {
	var tags []*model.Tag

	cmd := `
	  SELECT tags.id, tags.name 
		FROM tags
		  INNER JOIN recruitment_tags
			  ON tags.id = recruitment_tags.tag_id
		WHERE recruitment_tags.recruitment_id = $1
	`

	rows, err := dbConnection.Query(cmd, recId)
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
		tags = append(tags, &tag)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return tags, nil
}

func (t *Tag) CreateTag(dbConnection *sql.DB) (*model.Tag, error) {
	lower := strings.ToLower(t.Name)
	timeNow := time.Now().Local()

	cmd := "INSERT INTO tags (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, name"
	row := dbConnection.QueryRow(cmd, xid.New().String(), lower, timeNow, timeNow)

	var tag model.Tag
	err := row.Scan(&tag.ID, &tag.Name)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &tag, nil
}
