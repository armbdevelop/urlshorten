package middelware

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type responseWriterWrapper struct {
    http.ResponseWriter // оригинал
    statusCode     int
}

func (w *responseWriterWrapper) WriteHeader(code int) {
    w.statusCode = code          // сохраняем статус себе
    w.ResponseWriter.WriteHeader(code) // вызываем оригинальный метод
}


func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tn := time.Now()
		wrapper := &responseWriterWrapper{
            ResponseWriter: w,
            statusCode:     200, // дефолтный статус
        }
        next.ServeHTTP(wrapper, r)

		if wrapper.statusCode >= 500 {
			logger.Error("request failed",
				"status", wrapper.statusCode,
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(tn),
			)
		} else {
			logger.Info("request processed",
				"status", wrapper.statusCode,
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(tn),
			)
		}
	})
}