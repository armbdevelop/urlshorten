package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/armbdevelop/urlshorten/internal/service"
	"github.com/armbdevelop/urlshorten/internal/storage/memory"
)

func TestShortenerHandler_Shorten(t *testing.T) {
	// Arrange
	repo := memory.NewMemoryStorage()
	svc := service.NewShortenerService(repo)
	h := NewShortHandler(svc)

	body := `{"url":"https://google.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader([]byte(body)))
	rr := httptest.NewRecorder()

	// Act
	h.ShortenByOriginal(rr, req)

	// Assert
	if rr.Code != http.StatusCreated {
		t.Errorf("want %d, got %d", http.StatusCreated, rr.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(resp["short"]) != 10 {
		t.Errorf("want len 10, got %d", len(resp["short"]))
	}
}

func TestShorten_InvalidJSON(t *testing.T) {
	// Arrange
	repo := memory.NewMemoryStorage()
	svc := service.NewShortenerService(repo)
	h := NewShortHandler(svc)

	body := `{"url":"https://google.co`
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader([]byte(body)))
	rr := httptest.NewRecorder()

	// Act
	h.ShortenByOriginal(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("want %d, got %d", http.StatusCreated, rr.Code)
	}
}

func TestShorten_EmptyURL(t *testing.T) {
	// Arrange
	repo := memory.NewMemoryStorage()
	svc := service.NewShortenerService(repo)
	h := NewShortHandler(svc)

	body := ``
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader([]byte(body)))
	rr := httptest.NewRecorder()

	// Act
	h.ShortenByOriginal(rr, req)

	// Assert
	if rr.Code != http.StatusBadRequest {
		t.Errorf("unexpected error: %v, %v", http.StatusCreated, rr.Code)
	}
}

func TestGetOriginal_Success(t *testing.T) {
	// Arrange
	repo := memory.NewMemoryStorage()
	svc := service.NewShortenerService(repo)
	h := NewShortHandler(svc)

	// 1. Шаг создания ссылки (Вызываем обработчик создания)
	body := `{"url":"https://google.com"}`
	reqCreate := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader([]byte(body)))
	rrCreate := httptest.NewRecorder()

	// Вызываем метод создания
	h.ShortenByOriginal(rrCreate, reqCreate)

	if rrCreate.Code != http.StatusCreated && rrCreate.Code != http.StatusOK {
		t.Fatalf("failed to create short url, status: %d", rrCreate.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(rrCreate.Body).Decode(&resp); err != nil {
		t.Fatalf("decode create response: %v", err)
	}

	shortUrl := resp["short"]

	// 2. Шаг получения оригинальной ссылки
	reqGet := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/%v", shortUrl), nil)
	reqGet.SetPathValue("short", shortUrl)
	rrGet := httptest.NewRecorder()

	// Act - вызываем обработчик получения
	h.OriginalByShort(rrGet, reqGet)

	// Assert
	if rrGet.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, rrGet.Code)
	}

	var orig map[string]string
	if err := json.NewDecoder(rrGet.Body).Decode(&orig); err != nil {
		t.Fatalf("decode get response: %v", err)
	}

	if orig["original"] != "https://google.com" {
		t.Errorf("want url https://google.com, got %v", orig["original"])
	}
}

func TestGetOriginal_NotFound(t *testing.T) {
	// Arrange
	repo := memory.NewMemoryStorage()
	svc := service.NewShortenerService(repo)
	h := NewShortHandler(svc)

	body := `{"url":"https://google.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader([]byte(body)))
	rr := httptest.NewRecorder()

	// Act
	h.ShortenByOriginal(rr, req)

	reqGet := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/%v", "notexists"), nil)
	reqGet.SetPathValue("short", "notexists")
	rrGet := httptest.NewRecorder()

	// Act - вызываем обработчик получения
	h.OriginalByShort(rrGet, reqGet)

	// Assert
	if rrGet.Code != http.StatusNotFound {
		t.Errorf("want %d, got %d", http.StatusOK, rrGet.Code)
	}

}
