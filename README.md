<h1 align="center">Todo App</h1>

# Requirements

- `git`
- `docker`
- `docker-compose` (or Docker with Compose V2 support via `docker compose`)
- `make` (optional)
- `go` (optional)

# Getting Started

1. **Clone the repository**:

   ```
   $ git clone https://github.com/RX90/Todo-App.git
   $ cd Todo-App
   ```

2. **Build and start the application**:

   ```
   $ make build
   ```

3. <strong>Open in browser</strong>: <a href="http://localhost:8000" target="_blank">http://localhost:8000</a>

# Makefile Commands

| Command      | Description                                   |
| ------------ | --------------------------------------------- |
| `make build` | Build image and start containers              |
| `make up`    | Start containers without rebuilding           |
| `make down`  | Stop and remove containers with named volumes |
| `make start` | Start already built containers                |
| `make stop`  | Stop running containers                       |
| `make test`  | Run unit tests (need go)                      |

# Working Tree

    Todo-App
    ├─ .env
    ├─ client
    │  ├─ src
    │  │  ├─ font
    │  │  │  ├─ papyrus-pixel.ttf
    │  │  │  ├─ Rubik-Medium.ttf
    │  │  │  ├─ Rubik-Regular.ttf
    │  │  │  └─ Undertale-Battle-Font.ttf
    │  │  └─ img
    │  │     ├─ !done.svg
    │  │     ├─ black-delete.svg
    │  │     ├─ black-edit.svg
    │  │     ├─ checkbox.svg
    │  │     ├─ delete.svg
    │  │     ├─ done.svg
    │  │     ├─ dots.svg
    │  │     ├─ edit.svg
    │  │     ├─ hide.svg
    │  │     ├─ list.svg
    │  │     ├─ logo.ico
    │  │     ├─ plus.svg
    │  │     ├─ red-delete.svg
    │  │     ├─ search.svg
    │  │     ├─ sorry.jpg
    │  │     ├─ ugly-ico.ico
    │  │     ├─ ugly-logo.png
    │  │     ├─ ugly_logo_negate.png
    │  │     ├─ view.svg
    │  │     ├─ violet-checkbox.svg
    │  │     ├─ violet-delete.svg
    │  │     └─ violet-edit.svg
    │  ├─ static
    │  │  ├─ fetches.js
    │  │  ├─ main.css
    │  │  └─ main.js
    │  └─ templates
    │     └─ main.html
    ├─ docker-compose.yaml
    ├─ Dockerfile
    ├─ go.mod
    ├─ go.sum
    ├─ Makefile
    ├─ README.md
    ├─ server
    │  ├─ cmd
    │  │  └─ main.go
    │  ├─ configs
    │  │  └─ config.yaml
    │  ├─ internal
    │  │  ├─ app
    │  │  │  └─ app.go
    │  │  ├─ db
    │  │  │  └─ postgres.go
    │  │  ├─ handler
    │  │  │  ├─ auth.go
    │  │  │  ├─ auth_test.go
    │  │  │  ├─ handler.go
    │  │  │  ├─ list.go
    │  │  │  ├─ list_test.go
    │  │  │  ├─ middleware.go
    │  │  │  ├─ middleware_test.go
    │  │  │  ├─ task.go
    │  │  │  └─ task_test.go
    │  │  ├─ repository
    │  │  │  ├─ auth.go
    │  │  │  ├─ auth_test.go
    │  │  │  ├─ list.go
    │  │  │  ├─ list_test.go
    │  │  │  ├─ repository.go
    │  │  │  ├─ task.go
    │  │  │  └─ task_test.go
    │  │  ├─ service
    │  │  │  ├─ auth.go
    │  │  │  ├─ auth_test.go
    │  │  │  ├─ list.go
    │  │  │  ├─ mocks
    │  │  │  │  └─ mock.go
    │  │  │  ├─ service.go
    │  │  │  └─ task.go
    │  │  └─ todo
    │  │     └─ todo.go
    │  ├─ migrations
    │  │  ├─ 000001_init.down.sql
    │  │  └─ 000001_init.up.sql
    │  └─ server.go
    └─ wait-for-postgres.sh

Made by **[RX90](https://github.com/RX90)** && **[Mafiozich](https://github.com/Mafiozich)**
