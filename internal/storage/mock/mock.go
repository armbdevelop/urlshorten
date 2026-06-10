package mock

import (
	"context"

	"github.com/armbdevelop/urlshorten/internal/model"
)

// Storage — мок для тестов. Просто структура с функциями которые можно подменить.
type Storage struct {
	SaveFunc          func(ctx context.Context, url model.URL) error
	GetByShortFunc    func(ctx context.Context, short string) (model.URL, error)
	GetByOriginalFunc func(ctx context.Context, original string) (model.URL, error)
}

func (m *Storage) Save(ctx context.Context, url model.URL) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, url)
	}
	return nil
}

func (m *Storage) GetByShort(ctx context.Context, short string) (model.URL, error) {
	if m.GetByShortFunc != nil {
		return m.GetByShortFunc(ctx, short)
	}
	return model.URL{}, nil
}

func (m *Storage) GetByOriginal(ctx context.Context, original string) (model.URL, error) {
	if m.GetByOriginalFunc != nil {
		return m.GetByOriginalFunc(ctx, original)
	}
	return model.URL{}, nil
}
