FROM golang:1.24 AS builder
LABEL authors="Eugene"

# Изменим рабочую директорию на /app
WORKDIR /app/

# Копируем go.mod файлы для кеширования
COPY api-gateway/go.mod api-gateway/go.sum ./api-gateway/
COPY pkg/go.mod ./pkg/

# Копируем исходный код
COPY pkg/ ./pkg/
COPY api-gateway/ ./api-gateway/

# Переходим в директорию сервиса
WORKDIR /app/api-gateway

RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api-gateway ./cmd

FROM alpine:3.20

RUN apk add --no-cache curl ca-certificates

COPY --from=builder /app/api-gateway/api-gateway /usr/local/bin/api-gateway

RUN chmod +x /usr/local/bin/api-gateway

ENTRYPOINT ["api-gateway"]