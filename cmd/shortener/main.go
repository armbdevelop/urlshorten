package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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

	switch *storageType {
	case "memory":
		repo = memory.NewMemoryStorage()
	case "postgres":
		connStr := os.Getenv("DATABASE_URL")
		if connStr == "" {
			log.Fatal("DATABASE_URL is not set")
		}
		repo, err = postgres.NewPostgresStorage(connStr)
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown storage: %s", *storageType)
	}

	// 2. Сервис
	svc := service.NewShortenerService(repo)

	// 3. Хендлер
	h := handler.NewShortHandler(svc)

	// 4. Роутер
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/shorten", h.Shorten)
	mux.HandleFunc("GET /{short}", h.Redirect)

	appHandler := middelware.RequestID(middelware.Logger(middelware.Recovery(mux)))

	// 5. Сервер
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      appHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("Сервер запущен на http://localhost:8080")
	
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
