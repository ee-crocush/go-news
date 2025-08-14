# go-comments

Микросервис для управления комментариями в системе новостей, реализованный по принципам Domain-Driven Design (DDD).

## Назначение

Сервис обеспечивает полный жизненный цикл комментариев:
- Создание новых комментариев к новостям
- Получение комментариев по ID новости
- Модерацию комментариев через интеграцию с внешним сервисом модерации
- Управление статусами комментариев (ожидание, одобрено, отклонено)

## Структура проекта

```
├── Dockerfile                      # Docker образ для контейнеризации
├── cmd/
│   └── main.go                     # Точка входа в приложение
├── configs/
│   └── config.yaml                 # Конфигурационный файл
├── go.mod                          # Go модули
├── go.sum
├── internal/                       # Внутренняя логика приложения
│   ├── app/
│   │   └── run.go                  # Инициализация и запуск приложения
│   ├── domain/                     # Доменный слой (DDD)
│   │   └── comment/                # Агрегат комментариев
│   │       ├── comment.go          # Доменная модель комментария
│   │       ├── comment_test.go     # Тесты доменной модели
│   │       ├── contract.go         # Контракты и интерфейсы
│   │       ├── errors.go           # Доменные ошибки
│   │       ├── repository.go       # Интерфейс репозитория
│   │       ├── vo.go               # Value Objects
│   │       └── vo_test.go          # Тесты Value Objects
│   ├── infrastructure/             # Инфраструктурный слой
│   │   ├── config/
│   │   │   └── config.go           # Работа с конфигурацией
│   │   ├── events/                 # Kafka интеграция
│   │   │   └── comment.go          # События комментариев
│   │   ├── repo/                   # Реализации репозиториев
│   │   │   └── postgres/           # PostgreSQL репозиторий
│   │   │       ├── comment.go      # Реализация репозитория
│   │   │       ├── init.go         # Инициализация БД
│   │   │       └── mapper/         # Маппинг данных
│   │   │           └── comment.go  # Маппер для комментариев
│   │   └── transport/              # Транспортный слой
│   │       └── httplib/            # HTTP транспорт
│   │           ├── handler/        # HTTP обработчики
│   │           │   ├── create.go   # Создание комментария
│   │           │   ├── find_all_by_news_id.go # Получение по ID новости
│   │           │   ├── handler.go  # Базовый обработчик
│   │           │   └── health.go   # Health check
│   │           └── router.go       # Настройка маршрутизации
│   └── usecase/                    # Слой бизнес-логики (Use Cases)
│       └── comment/                # Use Cases для комментариев
│           ├── change_status.go    # Изменение статуса комментария
│           ├── create.go           # Создание комментария
│           ├── dto.go              # Data Transfer Objects
│           ├── find_all_by_news_id.go # Поиск по ID новости
│           └── interfaces.go       # Интерфейсы Use Cases
└── schema.sql                      # Схема базы данных
```

## Технологии

- **Go 1.21+** - основной язык разработки
- **PostgreSQL** - основное хранилище данных
- **Apache Kafka** - асинхронная интеграция с сервисом модерации
- **HTTP/REST** - API для взаимодействия с клиентами

## Локальная разработка

### Требования
- Go 1.21+
- PostgreSQL 14+
- Apache Kafka 2.8+

### Запуск приложения
```bash
# Установка зависимостей
go mod download

# Запуск сервиса
go run cmd/main.go
```

### Docker
```bash
# Сборка образа
docker build -t go-comments .

# Запуск контейнера
docker run -p 8081:8081 go-comments
```

## Конфигурация

Основная конфигурация находится в файле `configs/config.yaml` и `.env.example`.

## API Endpoints

### Комментарии
- `GET /comments/news/{id}` - получение всех комментариев для новости
- `POST /comments` - создание нового комментария

### Служебные
- `GET /health` - проверка состояния сервиса

### Примеры запросов

#### Получение комментариев
```bash
curl -X GET "http://localhost:8081/comments/news/1"
```

**Ответ:**
```json
{
  "comments": [
    {
      "id": "comment-uuid-456",
      "news_id": 1,
      "parent_id": null,
      "username": "username",
      "content": "Отличная новость!",
      "pub_time": "2024-01-01T10:00:00Z",
      "children": []
    }
  ]
}
```

#### Создание комментария
```bash
curl -X POST "http://localhost:8080/api/v1/comments" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Очень интересная статья",
    "news_id": "news-uuid-123"
  }'
```

**Ответ:**
```json
{
  "message": "Comment created successfully"
}
```

## Интеграции

### Kafka Publisher
Публикует события о созданных комментариях в топик модерации:

```json
{
  "comment_id": "comment-uuid-456",
  "content": "текст комментария", 
  "created_at": "2024-01-01T10:00:00Z"
}
```

### Kafka Consumer
Обрабатывает ответы от сервиса модерации:

```json
{
  "comment_id": "1",
  "status": "approved",
  "processed_at": "2024-01-01T10:00:00Z"
}
```

## Архитектура

Сервис построен по принципам Domain-Driven Design (DDD) и Clean Architecture:

- **Domain Layer** - доменные модели, value objects, бизнес-правила
- **Use Case Layer** - сценарии использования приложения
- **Infrastructure Layer** - внешние зависимости (БД, Kafka, HTTP)
- **Transport Layer** - входные точки (HTTP handlers)

### Принципы DDD
- **Агрегат Comment** - инкапсулирует логику работы с комментариями
- **Value Objects** - неизменяемые объекты (CommentContent, NewsID)
- **Repository Pattern** - абстракция для работы с хранилищем
- **Domain Events** - события для интеграции между bounded contexts

## Roadmap

### ✅ Реализовано
- Базовая архитектура DDD с разделением на слои
- Kafka интеграция для модерации
- HTTP API с REST endpoints
- PostgreSQL репозиторий с маппингом

### 🚧 Запланированные улучшения

#### Улучшение архитектуры DDD
- [ ] Закончить CRUD операции для комментариев
- [ ] Рефакторинг доменных агрегатов и усиление инвариантов

#### Покрытие тестами
- [ ] Unit тесты для доменной модели и value objects
- [ ] Integration тесты для репозиториев и БД
- [ ] Тесты для Kafka producer/consumer
- [ ] Мокирование внешних зависимостей

#### Система миграций БД
- [ ] Интеграция с golang-migrate для управления миграциями
- [ ] Версионирование схемы базы данных
- [ ] Автоматический rollback при ошибках миграции

#### Дополнительные возможности
- [ ] Метрики и мониторинг (Prometheus, Grafana)
- [ ] Улучшение логирования
- [ ] Комплексные health checks для всех зависимостей
- [ ] Rate limiting для защиты от злоупотреблений
- [ ] Кэширование популярных комментариев (Redis)
- [ ] Пагинация для больших списков комментариев