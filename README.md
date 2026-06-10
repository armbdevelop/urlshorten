# URL Shortener

Тестовое задание: сервис для создания коротких ссылок.

## Требования

- Go 1.25+
- Docker + Docker Compose (опционально, для PostgreSQL)
- Make (опционально)

## Технологии

- Go 1.25
- `net/http` (роутинг Go 1.22+)
- PostgreSQL 16 + `pgx/v5`
- In-memory storage (`sync.RWMutex`)
- Docker, Docker Compose
- Unit-tests (`testing`, `httptest`)

## API

### Создать короткую ссылку

```bash
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://google.com"}'
```

Ответ:
```json
{"short": "abc123_def"}
```

### Получить оригинальный URL

```bash
curl http://localhost:8080/api/abc123_def
```

Ответ:
```json
{"original": "https://google.com"}
```

### Идемпотентность

Повторный `POST` с тем же `url` вернёт ту же короткую ссылку.

## Переменные окружения

Скопируйте шаблон:

```bash
cp .env.example .env
```

Основные переменные:

| Переменная | Описание |
|------------|----------|
| `DATABASE_URL` | Строка подключения к PostgreSQL |
| `POSTGRES_USER` | Пользователь PostgreSQL |
| `POSTGRES_PASSWORD` | Пароль PostgreSQL |
| `POSTGRES_DB` | Имя базы данных |

## Запуск

### 1. Локально с in-memory хранилищем

```bash
go run ./cmd/shortener -storage=memory
```

Сервер будет доступен на `http://localhost:8080`.

### 2. Локально с PostgreSQL

```bash
# 1. Поднять PostgreSQL
docker-compose up -d

# 2. Запустить приложение
go run ./cmd/shortener -storage=postgres
```

### 3. Docker (in-memory)

```bash
docker build -t urlshortener .
docker run -p 8080:8080 urlshortener
```

### 4. Docker (PostgreSQL)

```bash
# 1. Поднять PostgreSQL
docker-compose up -d

# 2. Запустить приложение в Docker с нужным DATABASE_URL
docker run -p 8080:8080 \
  -e DATABASE_URL=postgres://urluser:urlpass@host.docker.internal:5432/urlshorten?sslmode=disable \
  urlshortener -storage=postgres
```

## Тесты

```bash
# Все тесты
go test ./...

# С подробным выводом
go test -v ./...

# С проверкой race condition и покрытием
go test -race -cover ./...
```

## Структура проекта

```
.
├── cmd/shortener/           # Точка входа
├── internal/
│   ├── handler/             # HTTP handlers
│   ├── middleware/          # Logging, recovery, request ID
│   ├── model/               # Структуры и ошибки
│   ├── service/             # Бизнес-логика + генерация short URL
│   └── storage/
│       ├── memory/          # In-memory реализация
│       └── postgres/        # PostgreSQL реализация
├── migrations/              # SQL миграции
├── Dockerfile
├── docker-compose.yml
├── .env.example
└── README.md
```

## Особенности реализации

- **Уникальность**: один оригинальный URL → одна короткая ссылка.
- **Генерация**: случайная 10-символьная строка из алфавита `a-zA-Z0-9_` через `crypto/rand`.
- **Конкурентность**: in-memory storage защищён `sync.RWMutex`.
- **Переключение хранилища**: через флаг `-storage=memory|postgres`.
