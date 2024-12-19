default:
	docker start todo-db
	cd server && go run cmd/main.go || exit 0
up:
	cd server && migrate -path ./internal/db/migrations -database "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" up
down:
	cd server && migrate -path ./internal/db/migrations -database "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" down
container:
	docker run --name=todo-db -e POSTGRES_PASSWORD='password' -p 5432:5432 -d postgres