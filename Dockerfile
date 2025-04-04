FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app ./server/cmd/main.go

FROM alpine:latest

RUN apk add --no-cache postgresql-client bash

WORKDIR /app

COPY --from=builder /app .

RUN chmod +x wait-for-postgres.sh

CMD ["./todo-app"]