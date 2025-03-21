<h1 align="center">Todo App</h1>

# Requirements:

### • [Docker](https://www.docker.com/products/docker-desktop/)

### • [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation)

### • [make](https://www.gnu.org/software/make/#download)

# Launching:

### First launch:

```
$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
$ export PATH=$PATH:$(go env GOPATH)/bin

$ make build
$ make run
```

#### After this you can go to 
```
localhost:8000
```

### Re-launch:

```
$ make run
```
Made by [RX90](https://github.com/RX90) and [Mafiozich](https://github.com/Mafiozich)
