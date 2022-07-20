package loader

import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nagokos/connefut_backend/graph/model"
	"github.com/nagokos/connefut_backend/graph/models/tag"
	"github.com/nagokos/connefut_backend/logger"
)

type TagReader struct {
	dbPool *pgxpool.Pool
}

func (u *TagReader) GetRecruitmentTags(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
	IDs := make([]interface{}, len(keys))
	var cmdArray []string
	for ix, key := range keys {
		IDs[ix] = key.String()
		cmdArray = append(cmdArray, fmt.Sprintf("$%d", ix+1))
	}

	tagByRecruitmentID, _ := tag.GetTagsByRecruitmentIDs(ctx, u.dbPool, IDs, cmdArray)

	output := make([]*dataloader.Result, len(keys))
	for index, recruitmentKey := range keys {
		output[index] = &dataloader.Result{Data: tagByRecruitmentID[recruitmentKey.String()], Error: nil}
	}
	return output
}

func LoadTagsByRecruitmentID(ctx context.Context, recruitmentID int) ([]*model.Tag, error) {
	loaders := GetLoaders(ctx)
	loaders.TagLoader.Clear(ctx, dataloader.StringKey(fmt.Sprintf("%d", recruitmentID))) // タグの更新ができなくなるため必要 募集のIDでキャッシュしているため同じIDでフィールドの値が変わるものは削除する
	thunk := loaders.TagLoader.Load(ctx, dataloader.StringKey(fmt.Sprintf("%d", recruitmentID)))
	result, err := thunk()
	if err != nil {
		logger.NewLogger().Error(err.Error())
		return nil, err
	}
	return result.([]*model.Tag), nil
}
