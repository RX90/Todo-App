# Todo App

**Для использования Todo-App надо:**

**1. Поднять базу данных на PostgreSQL с помощью Docker**

```
docker run --name=todo-db -e POSTGRES_PASSWORD='qwerty' -p 5432:5432 -d postgres - создаём контейнер

docker start todo-db - запускаем контейнер

cd server - переходим в директорию сервера

migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable" up - применяем 000001_init.up.sql
```

**2. Находясь в основной директории Todo-App, используем команды make для запуска Todo-App, make up для применения миграций и make down для отката миграций**

**3. Заходим в браузер и используем Todo-App по адресу localhost:8000**