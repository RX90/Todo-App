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

3. **Open in browser**: [http://localhost:8000](http://localhost:8000)

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

```
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
└─ server
   ├─ cmd
   │  └─ main.go
   ├─ configs
   │  └─ config.yaml
   ├─ internal
   │  ├─ app
   │  │  └─ app.go
   │  ├─ db
   │  │  └─ sqlite.go
   │  ├─ handler
   │  │  ├─ auth.go
   │  │  ├─ auth_test.go
   │  │  ├─ handler.go
   │  │  ├─ list.go
   │  │  ├─ list_test.go
   │  │  ├─ middleware.go
   │  │  ├─ middleware_test.go
   │  │  ├─ task.go
   │  │  └─ task_test.go
   │  ├─ repository
   │  │  ├─ auth.go
   │  │  ├─ auth_test.go
   │  │  ├─ list.go
   │  │  ├─ list_test.go
   │  │  ├─ repository.go
   │  │  ├─ task.go
   │  │  └─ task_test.go
   │  ├─ service
   │  │  ├─ auth.go
   │  │  ├─ auth_test.go
   │  │  ├─ list.go
   │  │  ├─ mocks
   │  │  │  └─ mock.go
   │  │  ├─ service.go
   │  │  └─ task.go
   │  └─ todo
   │     └─ todo.go
   ├─ migrations
   │  ├─ 000001_create_users_table.down.sql
   │  ├─ 000001_create_users_table.up.sql
   │  ├─ 000002_create_lists_table.down.sql
   │  ├─ 000002_create_lists_table.up.sql
   │  ├─ 000003_create_users_lists_table.down.sql
   │  ├─ 000003_create_users_lists_table.up.sql
   │  ├─ 000004_create_tasks_table.down.sql
   │  ├─ 000004_create_tasks_table.up.sql
   │  ├─ 000005_create_lists_tasks_table.down.sql
   │  ├─ 000005_create_lists_tasks_table.up.sql
   │  ├─ 000006_create_tokens_table.down.sql
   │  ├─ 000006_create_tokens_table.up.sql
   │  ├─ 000007_create_users_tokens_table.down.sql
   │  └─ 000007_create_users_tokens_table.up.sql
   └─ server.go
```

Made by **[RX90](https://github.com/RX90)** && **[Mafiozich](https://github.com/Cho-Nah)**
