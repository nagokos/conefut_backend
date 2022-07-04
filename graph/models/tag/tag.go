package tag

import (
	"context"
	"errors"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/utils"
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
		dbPool := db.DatabaseConnection()

		cmd := "SELECT COUNT(DISTINCT id) FROM tags WHERE name = $1"
		row := dbPool.QueryRow(context.Background(), cmd, lower)

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

func GetTags(ctx context.Context, dbPool *pgxpool.Pool) ([]*model.Tag, error) {
	var tags []*model.Tag

	cmd := "SELECT id, name FROM tags"
	rows, err := dbPool.Query(ctx, cmd)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag model.Tag
		err := rows.Scan(&tag.DatabaseID, &tag.Name)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		tag.ID = utils.GenerateAndSetUniqueID("Tag", *tag.DatabaseID)
		tags = append(tags, &tag)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return tags, nil
}

func (t *Tag) CreateTag(ctx context.Context, dbPool *pgxpool.Pool) (*model.Tag, error) {
	lower := strings.ToLower(t.Name)
	timeNow := time.Now().Local()

	cmd := "INSERT INTO tags (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id, name"
	row := dbPool.QueryRow(ctx, cmd, xid.New().String(), lower, timeNow, timeNow)

	var tag model.Tag
	err := row.Scan(&tag.ID, &tag.Name)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &tag, nil
}
