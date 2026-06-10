package service

import (
	"context"
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
	found, err := s.repo.GetByOriginal(ctx, originalURL)

	if err == nil {
		return found.ShortURL, nil
	}

	if err != model.ErrNotFound {
		return "", err
	}

	shortUrl := ""

	for {
		short, err := generateShort()
		if err != nil {
			return "", err
		}

		_, err = s.repo.GetByShort(ctx, short)
		if err == model.ErrNotFound {
			shortUrl = short
			break
		}
		if err != nil {
			return "", err
		}
		// коллизия — пробуем заново, вероятность крошечная но ненулевая
	}

	url := model.URL{
		OriginalURL: originalURL,
		ShortURL:    shortUrl,
		CreatedAt:   time.Now(),
	}

	err = s.repo.Save(ctx, url)

	if err != nil {
		return "", err
	}

	return url.ShortURL, nil

}

func (s *ShortenerService) GetOriginal(ctx context.Context, shortURL string) (string, error) {
	url, err := s.repo.GetByShort(ctx, shortURL)
	if err != nil {
		return "", err
	}
	return url.OriginalURL, nil
}
