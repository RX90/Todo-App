# Todo App

## Инструкция по использованию

### 1. Настройка базы данных PostgreSQL с помощью Docker

1. Создайте контейнер с базой данных:
   ```
   docker run --name=todo-db -e POSTGRES_PASSWORD='qwerty' -p 5432:5432 -d postgres
   ```
2. Запустите контейнер:
   ```
    docker start todo-db
   ```
3. Перейдите в директорию сервера:
   ```
   cd server
   ```
4. Примените миграции:
   ```
   migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable" up
   ```

### 2. Работа с приложением через Makefile

1. Запуск Todo-App:
   ```
   make
   ```
2. Применение миграций:
   ```
   make up
   ```
3. Откат миграций:
   ```
   make down
   ```

### 3. Использование приложения

Откройте браузер и перейдите по адресу: http://localhost:8000
