FROM golang:1.22-alpine AS builder

RUN apk add --no-cache build-base sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
ENV GOOS=linux

RUN go build -o todo-app ./server/cmd/main.go

FROM alpine:latest

RUN apk add --no-cache sqlite bash

WORKDIR /app

COPY --from=builder /app .

CMD ["./todo-app"]