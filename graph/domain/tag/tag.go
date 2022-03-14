package tag

import (
	"context"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/nagokos/connefut_backend/db"
	"github.com/nagokos/connefut_backend/ent"
	"github.com/nagokos/connefut_backend/ent/tag"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/logger"
)

type Tag struct {
	Name string
}

func existsTag() validation.RuleFunc {
	return func(v interface{}) error {
		var err error

		s := v.(string)
		ctx := context.Background()
		client := db.DatabaseConnection()

		res, _ := client.Tag.
			Query().
			Where(tag.NameEQ(s)).
			Exist(ctx)

		if res {
			err = errors.New("このタグは既に存在します")
		} else {
			err = nil
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

func GetTags(ctx context.Context, client *ent.Client) ([]*model.Tag, error) {
	var tags []*model.Tag

	res, err := client.Tag.
		Query().
		All(ctx)
	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("get tags error: %s", err.Error()))
		return nil, err
	}

	for _, tag := range res {
		tags = append(tags, &model.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		})
	}

	return tags, nil
}

func (t *Tag) CreateTag(ctx context.Context, client *ent.Client) (*model.Tag, error) {
	res, err := client.Tag.
		Create().
		SetName(t.Name).
		Save(ctx)

	if err != nil {
		logger.Log.Error().Msg(fmt.Sprintf("create tag error: %s", err.Error()))
		return nil, errors.New("タグの作成に失敗しました")
	}

	tag := &model.Tag{
		ID:   res.ID,
		Name: res.Name,
	}

	return tag, nil
}
