default:
	docker start todo-db
	cd server && go run cmd/main.go

migrate:
	migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable" up