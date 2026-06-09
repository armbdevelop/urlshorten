package middelware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

func generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
        return strconv.FormatInt(time.Now().UnixNano(), 36)
    }


	return hex.EncodeToString(b)
}

func RequestID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // до handler
		id := generateID()
		ctx := context.WithValue(r.Context(), "request_id", id)
		w.Header().Add("X-Request-ID", id)
        next.ServeHTTP(w, r.WithContext(ctx))
        // после handler
    })
}
