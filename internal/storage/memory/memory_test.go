package memory

import (
	"context"
	"testing"

	"github.com/armbdevelop/urlshorten/internal/model"
)

func TestMemoryStorage_SaveAndGet(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStorage()

	url := model.URL{
		OriginalURL: "https://google.com",
		ShortURL:    "abcdef123_",
	}

	err := s.Save(ctx, url)

	// Assert
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Act
	got, err := s.GetByShort(ctx, "abcdef123_")

	// Assert
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.OriginalURL != "https://google.com" {
		t.Errorf("want %q, got %q", "https://google.com", got.OriginalURL)
	}

}

func TestMemoryStorage_TestDuplicate(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStorage()

	url := model.URL{
		OriginalURL: "https://google.com",
		ShortURL:    "abcdef123_",
	}

	err := s.Save(ctx, url)

	// Assert
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// пытаемся заного сохранить чтобы словить ErrAlreadyExists
	err = s.Save(ctx, url)

	if err == nil {
		t.Fatalf("expected error for duplicate, got nil")
	}

	if err != model.ErrAlreadyExists {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}

}

func TestMemoryStorage_TestGetNotFound(t *testing.T) {
	ctx := context.Background()
	s := NewMemoryStorage()

	_, err := s.GetByOriginal(ctx, "sdfsdfs")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != model.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

}
