default:
	docker start todo-db
	cd server && go run cmd/main.go
up:
	cd server && migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable" up
down:
	cd server && migrate -path ./migrations -database "postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable" down