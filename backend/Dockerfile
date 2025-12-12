# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Установка зависимостей для миграций
RUN apk add --no-cache git

# Копирование go mod файлов
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/server ./cmd/server/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata postgresql-client

WORKDIR /app

# Установка goose
RUN apk add --no-cache go git && \
    go install github.com/pressly/goose/v3/cmd/goose@latest && \
    cp /root/go/bin/goose /usr/local/bin/goose && \
    apk del go git && \
    rm -rf /root/go

# Копирование бинарника
COPY --from=builder /app/bin/server .

# Копирование миграций
COPY --from=builder /app/migrations ./migrations

# Копирование entrypoint скрипта
COPY scripts/entrypoint.sh /app/entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]

