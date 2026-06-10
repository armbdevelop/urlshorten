package service

import (
	"context"
	"testing"

	"github.com/armbdevelop/urlshorten/internal/storage/memory"
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
