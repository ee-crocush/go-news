<h3>GoNews</h3>
Вариант решения практического задания по разработке агрегатора новостей.  

<h3>Установка mongoDB</h3>

Устанавливаем через docker compose. Пример compose:
```yaml
services:
  mongo_db:
    image: mongo
    container_name: mongo
    #    restart: always
    # environment:
    #   MONGO_INITDB_ROOT_USERNAME: ${MONGO_USER}
    #   MONGO_INITDB_ROOT_PASSWORD: ${MONGO_PASSWORD}
    volumes:
      - mongo_data:/data/db
    ports:
      - '27017:27017'

volumes:
  mongo_data:

```

Если захочется, можно добавить .env:
```dotenv
MONGO_USER=user
MONGO_PASSWORD=password
```