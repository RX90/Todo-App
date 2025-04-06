build:
	docker-compose up -d --build

up:
	docker-compose up -d

down:
	docker-compose down -v

start:
	docker-compose start

stop:
	docker-compose stop

test:
	go clean -testcache
	go test -v ./...