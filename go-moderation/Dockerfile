FROM golang:1.24 AS builder
LABEL authors="Eugene"

# Изменим рабочую директорию на /app
WORKDIR /app/

# Копируем go.mod файлы для кеширования
COPY go-moderation/go.mod go-moderation/go.sum ./go-moderation/
COPY pkg/go.mod ./pkg/

# Копируем исходный код
COPY pkg/ ./pkg/
COPY go-moderation/ ./go-moderation/

# Переходим в директорию сервиса
WORKDIR /app/go-moderation

RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-moderation ./cmd

FROM alpine:3.20

RUN apk add --no-cache curl ca-certificates

COPY --from=builder /app/go-moderation/go-moderation /usr/local/bin/go-moderation

RUN chmod +x /usr/local/bin/go-moderation

ENTRYPOINT ["go-moderation"]