package memory

import (
	"context"
	"sync"

	"github.com/armbdevelop/urlshorten/internal/model"
)

type MemoryStorage struct {
	mu         sync.RWMutex
	byShort    map[string]model.URL // short → URL
	byOriginal map[string]model.URL // original → URL
}


func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{
        byShort:    make(map[string]model.URL),
        byOriginal: make(map[string]model.URL),
    }
}


func (s *MemoryStorage) Save(ctx context.Context, url model.URL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.byOriginal[url.OriginalURL]

	if ok {
		return model.ErrAlreadyExists
	}
	
	_, ok = s.byShort[url.ShortURL]

	if ok {
		return model.ErrAlreadyExists
	}

    s.byOriginal[url.OriginalURL] = url
    s.byShort[url.ShortURL] = url

	return nil
}


func (s *MemoryStorage) GetByShort(ctx context.Context, short string) (model.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	found, ok := s.byShort[short]

	if !ok {
		return model.URL{}, model.ErrNotFound
	}

	return found, nil
}

func (s *MemoryStorage) GetByOriginal(ctx context.Context, original string) (model.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	found, ok := s.byOriginal[original]

	if !ok {
		return model.URL{}, model.ErrNotFound
	}

	return found, nil
}