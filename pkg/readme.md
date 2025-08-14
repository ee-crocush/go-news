# pkg

Общая библиотека компонентов для микросервисов новостной системы. Содержит переиспользуемые компоненты инфраструктуры, которые используются всеми сервисами системы.

## Назначение

Библиотека предоставляет стандартизированные компоненты для:
- Конфигурирование приложений с поддержкой переменных окружения
- Структурированное логирование с настраиваемыми уровнями
- Kafka интеграция (Consumer/Producer) с готовыми конфигурациями
- HTTP middleware для общих задач (CORS, логирование, метрики)
- Запуск HTTP серверов на базе Fiber с едиными настройками
- API утилиты для стандартизации ответов

## Структура проекта

```
├── api/                            # API утилиты и хелперы
│   └── api.go                      # Стандартные ответы API
├── config/                         # Загрузка и валидация конфигурации
│   └── loader.go                   # Загрузчик конфигов с env поддержкой
├── go.mod                          # Go модули
├── go.sum
├── kafka/                          # Kafka интеграция
│   ├── consumer.go                 # Kafka Consumer с автоконфигурацией
│   └── publisher.go                # Kafka Publisher с retry логикой
├── logger/                         # Структурированное логирование
│   └── logger.go                   # Настройка логгера (logrus/zap)
├── middleware/                     # HTTP middleware компоненты
│   └── middleware.go               # CORS, Request ID, Recovery, Logging
└── server/                         # HTTP серверы
    ├── fiber/
    │   └── fiber.go                # Fiber сервер с настройками
    └── server.go                   # Запуск и Graceful shutdown
```

## Roadmap

### ✅ Реализовано
- Базовые компоненты инфраструктуры для всех сервисов
- Kafka Producer/Consumer с автоконфигурацией
- HTTP сервер на Fiber с middleware и Graceful shutdown
- Структурированное логирование
- Загрузчик конфигурации с env поддержкой
- Стандартизированные API ответы

### 🚧 Возможные улучшения

#### Расширение middleware
- [ ] JWT аутентификация middleware
- [ ] Rate limiting middleware с Redis
- [ ] Circuit breaker для внешних сервисов
- [ ] Request/Response validation middleware
- [ ] Compression middleware (gzip, brotli)
- [ ] Security headers middleware (HSTS, CSP)

#### Мониторинг и метрики
- [ ] Prometheus метрики middleware
- [ ] Health check utilities
- [ ] Готовые дашборды Grafana
- [ ] Structured metrics для бизнес-событий

#### Database utilities
- [ ] MongoDB connection pool и helpers
- [ ] PostgreSQL connection pool и migration tools
- [ ] Redis client с connection pooling
- [ ] Database health checks
- [ ] Transaction helpers и patterns

#### Kafka улучшения
- [ ] Schema Registry интеграция
- [ ] Kafka Admin API для управления топиками

#### Testing utilities
- [ ] Test containers для интеграционных тестов
- [ ] Mock generators для интерфейсов
- [ ] HTTP test helpers
- [ ] Kafka test utilities
- [ ] Database test fixtures

#### Безопасность
- [ ] Encryption/Decryption utilities
- [ ] Password hashing helpers (bcrypt, argon2)
- [ ] API key validation
- [ ] OAuth2/OIDC интеграция
- [ ] Input sanitization helpers

#### Дополнительные утилиты
- [ ] Configuration hot reload

#### Документация и примеры
- [ ] Подробная документация по всем компонентам
- [ ] Примеры использования для каждого сервиса
- [ ] Best practices руководство
- [ ] Migration guide между версиями
- [ ] API документация с примерами