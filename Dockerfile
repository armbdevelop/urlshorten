# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Сначала go.mod/go.sum для кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем
RUN go build -o urlshortener ./cmd/shortener

# Final stage
FROM alpine:latest

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/urlshortener .

# Порт
EXPOSE 8080

# Дефолт: memory. Переопределяй через command в docker-compose
ENTRYPOINT ["./urlshortener"]
CMD ["-storage=memory"]
