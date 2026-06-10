package handler

import (
	"encoding/json"
	"net/http"

	"github.com/armbdevelop/urlshorten/internal/service"
)

type ShortenerHandler struct {
	service *service.ShortenerService
}

func NewShortHandler(service *service.ShortenerService) *ShortenerHandler {
	return &ShortenerHandler{service: service}
}

func (h *ShortenerHandler) ShortenByOriginal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var url shortenRequest

	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	if url.URL == "" {
		http.Error(w, "url is required", http.StatusBadRequest)
		return
	}

	shortUrl, err := h.service.Shorten(ctx, url.URL)

	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"short": shortUrl,
	})

}

func (h *ShortenerHandler) OriginalByShort(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	short := r.PathValue("short")

	if short == "" {
		http.Error(w, "short is required", http.StatusBadRequest)
		return
	}

	original, err := h.service.GetOriginal(ctx, short)

	if err != nil {
		handleServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"original": original,
	})

}
