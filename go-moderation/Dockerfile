FROM golang:1.24 AS builder
LABEL authors="Eugene"

WORKDIR /app/

COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

COPY ./go-moderation /app/
COPY ./pkg /app/

RUN go mod tidy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./


FROM alpine:3.20

RUN apk add --no-cache curl

COPY --from=builder /app/go-news /usr/local/bin/go-news

RUN chmod +x /usr/local/bin/go-news

ENTRYPOINT [ "go-news" ]