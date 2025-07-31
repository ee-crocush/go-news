# Frontend для микросервисов go-news

## Локальная разработка

### Установка зависимостей

```
npm install
```

### Запуск в dev режиме на локальном хосте

```
npm run dev
```

### Сборка исходников

```
npm run build
```

### Запуск линтеров

```
npm run lint
```

## Локальный запуск в Docker

1. Для запуска нужно удостовериться, что сеть `go-news_network` создана, иначе создаем с помощью команды 
`docker network create go-news_network`.
2. Создаем образ командой `docker compose build` либо `docker build -t <image-name>:<tag> <path-to-Dockerfile>`
3. Далее Фронтенд сервис можно запустить в docker командой `docker compose up -d`, либо если пропустить второй пункт 
запустить можно командой `docker compose up -d --build`
