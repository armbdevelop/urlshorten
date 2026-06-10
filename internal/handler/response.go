package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/armbdevelop/urlshorten/internal/model"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	body, err := json.Marshal(data)
	if err != nil {
		http.Error(w, `{"error":"encode error"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(body)
}

func writeErrorJSON(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrNotFound):
		writeErrorJSON(w, http.StatusNotFound, err.Error())
	case errors.Is(err, model.ErrAlreadyExists):
		writeErrorJSON(w, http.StatusConflict, err.Error())
	default:
		log.Printf("internal error: %v", err)
		writeErrorJSON(w, http.StatusInternalServerError, "internal server error")
	}
}
