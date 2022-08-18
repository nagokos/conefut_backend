package gcs

import (
	"context"

	"cloud.google.com/go/storage"
)

func NewGcsClient(ctx context.Context) (*storage.Client, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}
