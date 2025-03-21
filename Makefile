run:
	docker start todo-db
	cd server && go run cmd/main.go || exit 0
build:
	docker run --name=todo-db -e POSTGRES_PASSWORD='password' -p 5432:5432 -d postgres
	sleep 2
	migrate -path ./server/internal/db/migrations -database "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" up
