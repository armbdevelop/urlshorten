package service

import (
	"context"
	"testing"

	"github.com/armbdevelop/urlshorten/internal/model"
	"github.com/armbdevelop/urlshorten/internal/storage/memory"
	"github.com/armbdevelop/urlshorten/internal/storage/mock"
)

func TestShorten_NewURL(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewMemoryStorage()
	svc := NewShortenerService(repo)

	short, err := svc.Shorten(ctx, "https://google.com")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(short) != 10 {
		t.Errorf("want len 10, got %d", len(short))
	}

}

func TestShorten_Duplicate(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewMemoryStorage()
	svc := NewShortenerService(repo)

	url, err := svc.Shorten(ctx, "https://google.com")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	url2, err := svc.Shorten(ctx, "https://google.com")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if url != url2 {
		t.Fatalf("unexpected error: want %v, got %v", url, url2)
	}

}

func TestGetOriginal_NotFound(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewMemoryStorage()
	svc := NewShortenerService(repo)

	svc.Shorten(ctx, "https://ya.ru")

	_, err := svc.GetOriginal(ctx, "NOTEXIST")

	if err == nil {
		t.Fatalf("expected error: %v", err)
	}

}

func TestGetOriginal_Found(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewMemoryStorage()
	svc := NewShortenerService(repo)

	url, err := svc.Shorten(ctx, "https://ya.ru")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = svc.GetOriginal(ctx, url)

	if err != nil {
		t.Fatalf("expected error: %v", err)
	}

}

// --- unit тесты с моком ---

func TestShorten_Unit_ExistingURL(t *testing.T) {
	ctx := context.Background()

	repo := &mock.Storage{
		GetByOriginalFunc: func(ctx context.Context, original string) (model.URL, error) {
			return model.URL{OriginalURL: "https://ya.ru", ShortURL: "abc123_def"}, nil
		},
	}

	svc := NewShortenerService(repo)
	short, err := svc.Shorten(ctx, "https://ya.ru")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if short != "abc123_def" {
		t.Errorf("want abc123_def, got %s", short)
	}
}

func TestShorten_Unit_SaveOK(t *testing.T) {
	ctx := context.Background()

	repo := &mock.Storage{
		GetByOriginalFunc: func(ctx context.Context, original string) (model.URL, error) {
			return model.URL{}, model.ErrNotFound
		},
		SaveFunc: func(ctx context.Context, url model.URL) error {
			return nil
		},
	}

	svc := NewShortenerService(repo)
	short, err := svc.Shorten(ctx, "https://new-url.com")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(short) != 10 {
		t.Errorf("want len 10, got %d", len(short))
	}
}

func TestShorten_Unit_RetryOnDuplicateShort(t *testing.T) {
	ctx := context.Background()
	callCount := 0

	repo := &mock.Storage{
		GetByOriginalFunc: func(ctx context.Context, original string) (model.URL, error) {
			return model.URL{}, model.ErrNotFound
		},
		SaveFunc: func(ctx context.Context, url model.URL) error {
			callCount++
			if callCount == 1 {
				return model.ErrAlreadyExists // коллизия short
			}
			return nil
		},
	}

	svc := NewShortenerService(repo)
	short, err := svc.Shorten(ctx, "https://retry.com")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(short) != 10 {
		t.Errorf("want len 10, got %d", len(short))
	}
	if callCount != 2 {
		t.Errorf("want 2 save calls, got %d", callCount)
	}
}
