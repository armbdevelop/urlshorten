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
        http.Error(w, "encode error", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _, _ = w.Write(body)
}


func handleServiceError(w http.ResponseWriter, err error) {
    switch {
    case errors.Is(err, model.ErrNotFound):
        http.Error(w, err.Error(), http.StatusNotFound)
    case errors.Is(err, model.ErrAlreadyExists):
        http.Error(w, err.Error(), http.StatusConflict)
    default:
        log.Printf("internal error: %v", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
    }
}