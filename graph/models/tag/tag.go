package tag

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/utils"
	"github.com/nagokos/connefut_backend/logger"
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

func GetTags(ctx context.Context, dbPool *pgxpool.Pool) (*model.TagConnection, error) {
	connection := model.TagConnection{
		PageInfo: &model.PageInfo{},
	}

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
		connection.Edges = append(connection.Edges, &model.TagEdge{
			Cursor: utils.GenerateUniqueID("Tag", tag.DatabaseID),
			Node:   &tag,
		})
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	return &connection, nil
}

func GetTagsByRecruitmentIDs(ctx context.Context, dbPool *pgxpool.Pool, IDs []interface{}, cmdArray []string) (map[string][]*model.Tag, error) {
	if len(IDs) == 0 {
		return nil, nil
	}

	cmd := fmt.Sprintf(
		`
	  SELECT t.id, t.name, r_t.recruitment_id
		FROM tags AS t
		INNER JOIN recruitment_tags AS r_t
		ON r_t.tag_id = t.id
		WHERE r_t.recruitment_id IN (%s)
		`,
		strings.Join(cmdArray, ","),
	)

	rows, err := dbPool.Query(
		ctx,
		cmd,
		IDs...,
	)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, nil
	}
	defer rows.Close()

	tagByRecruitmentID := map[string][]*model.Tag{}
	for rows.Next() {
		var tag model.Tag
		var recruitmentID int
		err := rows.Scan(&tag.DatabaseID, &tag.Name, &recruitmentID)
		if err != nil {
			logger.NewLogger().Error(err.Error())
		}
		tagByRecruitmentID[strconv.Itoa(recruitmentID)] = append(tagByRecruitmentID[strconv.Itoa(recruitmentID)], &tag)
	}

	err = rows.Err()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, nil
	}

	return tagByRecruitmentID, nil
}

func (t *Tag) CreateTag(ctx context.Context, dbPool *pgxpool.Pool) (*model.CreateTagPayload, error) {
	timeNow := time.Now().Local()

	cmd := "INSERT INTO tags (name, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id, name"
	row := dbPool.QueryRow(ctx, cmd, t.Name, timeNow, timeNow)

	var tag model.Tag
	err := row.Scan(&tag.DatabaseID, &tag.Name)
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}

	payload := model.CreateTagPayload{
		FeedbackTagEdge: &model.TagEdge{
			Cursor: utils.GenerateUniqueID("Tag", tag.DatabaseID),
			Node:   &tag,
		},
	}
	return &payload, nil
}
