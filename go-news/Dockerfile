FROM golang:1.24 AS builder
LABEL authors="Eugene"

# Изменим рабочую директорию на /app
WORKDIR /app/

# Копируем go.mod файлы для кеширования
COPY go-news/go.mod go-news/go.sum ./go-news/
COPY pkg/go.mod ./pkg/

# Копируем исходный код
COPY pkg/ ./pkg/
COPY go-news/ ./go-news/

# Переходим в директорию сервиса
WORKDIR /app/go-news

RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-news ./cmd

FROM alpine:3.20

RUN apk add --no-cache curl ca-certificates

COPY --from=builder /app/go-news/go-news /usr/local/bin/go-news

RUN chmod +x /usr/local/bin/go-news

ENTRYPOINT ["go-news"]