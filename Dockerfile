FROM golang:1.22

RUN apt-get update && apt-get -y install postgresql-client

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN chmod +x wait-for-postgres.sh
RUN go build -o todo-app ./server/cmd/main.go

CMD ["./todo-app"]