package model

import "time"

type URL struct {
	ID          int64
	OriginalURL string
	ShortURL    string
	CreatedAt   time.Time
}
