package storage

import (
	"context"

	"github.com/armbdevelop/urlshorten/internal/model"
)

type Repository interface {
	Save(ctx context.Context, url model.URL) error
	GetByShort(ctx context.Context, short string) (model.URL, error)
	GetByOriginal(ctx context.Context, original string) (model.URL, error)
}
