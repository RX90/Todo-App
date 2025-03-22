run:
	docker start todo-db
	go run server/cmd/main.go || exit 0
build:
	docker run --name=todo-db -e POSTGRES_PASSWORD='password' -p 5432:5432 -d postgres
	sleep 3
	migrate -path ./server/migrations -database "postgres://postgres:password@localhost:5432/postgres?sslmode=disable" up
test:
	go clean -testcache
	go test -v ./...