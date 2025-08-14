FROM golang:1.24 AS builder
LABEL authors="Eugene"

# Изменим рабочую директорию на /app
WORKDIR /app/

# Копируем go.mod файлы для кеширования
COPY go-comments/go.mod go-comments/go.sum ./go-comments/
COPY pkg/go.mod ./pkg/

# Копируем исходный код
COPY pkg/ ./pkg/
COPY go-comments/ ./go-comments/

# Переходим в директорию сервиса
WORKDIR /app/go-comments

RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-comments ./cmd

FROM alpine:3.20

RUN apk add --no-cache curl ca-certificates

COPY --from=builder /app/go-comments/go-comments /usr/local/bin/go-comments

RUN chmod +x /usr/local/bin/go-comments

ENTRYPOINT ["go-comments"]