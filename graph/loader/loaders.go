package loader

import (
	"context"
	"net/http"

	"github.com/graph-gophers/dataloader"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

type Loaders struct {
	UserLoader       *dataloader.Loader
	PrefectureLoader *dataloader.Loader
	SportLoader      *dataloader.Loader
	TagLoader        *dataloader.Loader
}

func NewLoaders(dbPool *pgxpool.Pool) *Loaders {
	userReader := &UserReader{dbPool: dbPool}
	prefectureReader := &PrefectureReader{dbPool: dbPool}
	sportReader := &SportReader{dbPool: dbPool}
	TagReader := &TagReader{dbPool: dbPool}
	loaders := &Loaders{
		UserLoader:       dataloader.NewBatchedLoader(userReader.GetUsers),
		PrefectureLoader: dataloader.NewBatchedLoader(prefectureReader.GetPrefectures),
		SportLoader:      dataloader.NewBatchedLoader(sportReader.GetSports),
		TagLoader:        dataloader.NewBatchedLoader(TagReader.GetRecruitmentTags),
	}
	return loaders
}

func Middleware(loaders *Loaders, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCtx := context.WithValue(r.Context(), loadersKey, loaders)
		r = r.WithContext(nextCtx)
		next.ServeHTTP(w, r)
	})
}

func GetLoaders(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
