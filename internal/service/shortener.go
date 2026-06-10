package service

import (
	"context"
	"errors"
	"time"

	"github.com/armbdevelop/urlshorten/internal/model"
	"github.com/armbdevelop/urlshorten/internal/storage"
)

type ShortenerService struct {
	repo storage.Repository
}

func NewShortenerService(repo storage.Repository) *ShortenerService {
	return &ShortenerService{repo: repo}
}

func (s *ShortenerService) Shorten(ctx context.Context, originalURL string) (string, error) {
	// может уже есть такой original
	found, err := s.repo.GetByOriginal(ctx, originalURL)
	if err == nil {
		return found.ShortURL, nil
	}
	if err != model.ErrNotFound {
		return "", err
	}

	for {
		short, err := generateShort()
		if err != nil {
			return "", err
		}

		url := model.URL{
			OriginalURL: originalURL,
			ShortURL:    short,
			CreatedAt:   time.Now(),
		}

		// сразу сохраняем, не проверяя заранее
		err = s.repo.Save(ctx, url)
		if err == nil {
			return short, nil
		}

		// если уже есть — проверяем, может это тот же original вставил другой поток
		if errors.Is(err, model.ErrAlreadyExists) {
			found, err = s.repo.GetByOriginal(ctx, originalURL)
			if err == nil {
				return found.ShortURL, nil
			}
			// не тот же original — значит коллизия short, пробуем заново
			continue
		}

		return "", err
	}
}

func (s *ShortenerService) GetOriginal(ctx context.Context, shortURL string) (string, error) {
	url, err := s.repo.GetByShort(ctx, shortURL)
	if err != nil {
		return "", err
	}
	return url.OriginalURL, nil
}
