# Новостной агрегатор

## Обзор проекта

Новостной агрегатор - это демонстрационное распределенное приложение, построенное на микросервисной архитектуре с 
применением принципов Domain-Driven Design (DDD). Система предназначена для сбора, обработки и предоставления 
новостного контента с возможностью комментирования.

## Архитектура системы

Приложение состоит из frontend-части и трех основных backend-сервисов, каждый из которых отвечает за определенную 
функциональную область

![схема.drawio.png](doc/%D1%81%D1%85%D0%B5%D0%BC%D0%B0.drawio.png)

### Frontend (Vue.js + Nginx)
**Назначение:** Пользовательский интерфейс приложения
- **Vue.js приложение:** SPA (Single Page Application) для отображения новостей и взаимодействия с пользователем
- **Nginx:** Веб-сервер для обслуживания статических файлов
- Интерактивный интерфейс для чтения новостей и комментирования
- Система фильтрации и поиска новостей (корректировка)

### API Gateway
**Назначение:** Единая точка входа для всех клиентских запросов
- Маршрутизация запросов к соответствующим микросервисам
- Агрегация ответов от различных сервисов
- Кэширование часто запрашиваемых данных
- Rate limiting и защита от DDoS-атак

### go-news
**Назначение:** Основной сервис управления новостным контентом
- Сбор новостей из различных источников (RSS)
- Обработка и нормализация новостных данных

### go-comments
**Назначение:** Сервис для работы с пользовательскими комментариями
- Создание комментариев
- Система модерации комментариев
- Поддержка вложенных комментариев (древовидная структура)

## Технический стек

**Backend:**
- **Язык программирования:** Go (Golang)
- **Архитектурный подход:** Domain-Driven Design (DDD)
- **Паттерн архитектуры:** Микросервисы
- **Коммуникация между сервисами:** HTTP REST API
- **База данных:** PostgreSQL, MongoDB
- **Очереди сообщений:** Apache Kafka

**Frontend:**
- **Фреймворк:** Vue.js 3
- **Веб-сервер:** Nginx
- **Сборщик:** Vite/Webpack
- **UI библиотеки:** Vuetify

**DevOps:**
- **Контейнеризация:** Docker
- **Оркестрация:** Kubernetes
- **Reverse Proxy:** Nginx (для статики и проксирования API запросов)