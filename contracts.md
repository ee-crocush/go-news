# API Контракты

Документация по REST API endpoints для системы новостей.

## Base URL
```
http://localhost:8080
```

## Endpoints

### 1. Получение списка новостей
```json
{
  "method": "GET",
  "url": "/api/news/",
  "parameters": {
    "page": {
      "type": "number",
      "required": false,
      "default": 1,
      "description": "Номер страницы"
    },
    "search": {
      "type": "string",
      "required": false,
      "description": "Поиск по заголовку новости"
    },
    "limit": {
      "type": "number",
      "required": false,
      "default": 10,
      "minimum": 10,
      "description": "Количество новостей на странице"
    }
  },
  "response": {
    "data": {
      "news": [
        {
          "id": "number",
          "title": "string",
          "content": "string",
          "link": "string",
          "pub_time": "string"
        }
      ],
      "total": "number"
    }
  }
}
```

### 2. Получение последней новости
```json
{
  "method": "GET",
  "url": "/api/news/last",
  "parameters": {},
  "response": {
    "data": {
      "post": {
        "id": "number",
        "title": "string",
        "content": "string",
        "link": "string",
        "pub_time": "string"
      }
    }
  }
}
```

### 3. Получение последних новостей
```json
{
  "method": "GET",
  "url": "/api/news/latest",
  "parameters": {
    "limit": {
      "type": "number",
      "required": false,
      "minimum": 10,
      "description": "Количество последних новостей"
    }
  },
  "response": {
    "data": {
      "news": [
        {
          "id": "number",
          "title": "string",
          "content": "string",
          "link": "string",
          "pub_time": "string"
        }
      ]
    }
  }
}
```

### 4. Получение детальной информации о новости
```json
{
  "method": "GET",
  "url": "/api/news/{id}",
  "parameters": {
    "id": {
      "type": "number",
      "required": true,
      "location": "path",
      "description": "ID новости"
    }
  },
  "response": {
    "data": {
      "post": {
        "id": "number",
        "title": "string",
        "content": "string",
        "link": "string",
        "pub_time": "string"
      },
      "comments": [
        {
          "id": "number",
          "news_id": "number",
          "parent_id": "number|null",
          "username": "string",
          "content": "string",
          "pub_time": "string",
          "children": [
            {
              "id": "number",
              "news_id": "number",
              "parent_id": "number",
              "username": "string",
              "content": "string",
              "pub_time": "string",
              "children": []
            }
          ]
        }
      ]
    }
  }
}
```

### 5. Создание комментария
```json
{
  "method": "POST",
  "url": "/api/comments",
  "headers": {
    "Content-Type": "application/json"
  },
  "body": {
    "news_id": "number (required)",
    "parent_id": "number (optional)",
    "username": "string (required)",
    "content": "string (required)"
  },
  "response": {
    "status": "string",
    "message": "string"
  }
}
```

## Примеры запросов

### Получение новостей с пагинацией
```bash
GET /api/news/?page=2&limit=20
```

### Поиск новостей
```bash
GET /api/news/?search=технологии&page=1&limit=15
```

### Получение новости с комментариями
```bash
GET /api/news/123
```

### Создание комментария
```bash
POST /api/comments
Content-Type: application/json

{
  "news_id": 123,
  "parent_id": 456,
  "username": "user123",
  "content": "Отличная новость!"
}
```

### Создание комментария верхнего уровня
```bash
POST /api/comments
Content-Type: application/json

{
  "news_id": 123,
  "username": "user123",
  "content": "Первый комментарий к новости"
}
```