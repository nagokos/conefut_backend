package tag

import (
	"database/sql"
	"errors"
	"fmt"
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
		_, dbConnection := db.DatabaseConnection()

		cmd := fmt.Sprintf("SELECT COUNT(DISTINCT id) FROM %s WHERE name = $1", db.TagTable)
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

	cmd := fmt.Sprintf("SELECT id, name FROM %s", db.TagTable)
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

	cmd := fmt.Sprintf(`
	  SELECT tags.id, tags.name 
		FROM %s
		  INNER JOIN %s
			  ON tags.id = recruitment_tags.tag_id
		WHERE recruitment_tags.recruitment_id = $1
	`, db.TagTable, db.RecruitmentTagsTable)

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

	cmd := fmt.Sprintf("INSERT INTO %s (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, name", db.TagTable)
	row := dbConnection.QueryRow(cmd, xid.New().String(), lower, timeNow, timeNow)

	var tag model.Tag
	err := row.Scan(&tag.ID, &tag.Name)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &tag, nil
}
