package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/armbdevelop/urlshorten/internal/handler"
	"github.com/armbdevelop/urlshorten/internal/middleware"
	"github.com/armbdevelop/urlshorten/internal/service"
	"github.com/armbdevelop/urlshorten/internal/storage"
	"github.com/armbdevelop/urlshorten/internal/storage/memory"
	"github.com/armbdevelop/urlshorten/internal/storage/postgres"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	// 1. Флаг -storage
	storageType := flag.String("storage", "memory", "storage type: memory or postgres")
	flag.Parse()

	var repo storage.Repository
	var err error
	var pgRepo *postgres.PostgresStorage

	switch *storageType {
	case "memory":
		repo = memory.NewMemoryStorage()
	case "postgres":
		connStr := os.Getenv("DATABASE_URL")
		if connStr == "" {
			log.Fatal("DATABASE_URL is not set")
		}
		pgRepo, err = postgres.NewPostgresStorage(connStr)
		if err != nil {
			log.Fatal(err)
		}
		repo = pgRepo
	default:
		log.Fatalf("unknown storage: %s", *storageType)
	}

	// 2. Сервис
	svc := service.NewShortenerService(repo)

	// 3. Хендлер
	h := handler.NewShortHandler(svc)

	// 4. Роутер
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/shorten", h.ShortenByOriginal)
	mux.HandleFunc("GET /api/{short}", h.OriginalByShort)

	appHandler := middelware.RequestID(middelware.Logger(middelware.Recovery(mux)))

	// 5. Сервер
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      appHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// ловим сигналы для graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Println("Сервер запущен на http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	// даем серверу 5 секунд на завершение текущих запросов
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	if pgRepo != nil {
		pgRepo.Close()
	}

	log.Println("server stopped")
}
