<h1 align="center">Todo App</h1>

<h2 align="center">Инструкция по настройке и использованию Todo App</h2>

### 0. Для корректной настройки нужно установить <a href="https://www.docker.com/products/docker-desktop/" target="_blank">Docker</a> и <a href="https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation" target="_blank">golang-migrate</a>

---

### 1. Настройка и запуск проекта

1. Создайте контейнер с базой данных:
   ```
   $ make container
   ```
2. Примените миграции:
   ```
   $ make up
   ```
3. Запустите Todo-App:
   ```
   make
   ```
4. Откройте браузер и перейдите по адресу http://localhost:8000

---

### 2. Работа с приложением через Makefile

1. Запуск Todo-App:
   ```
   $ make
   ```
2. Применение миграций:
   ```
   $ make up
   ```
3. Откат миграций:
   ```
   $ make down
   ```
4. Создание контейнера:
   ```
   make container
   ```
