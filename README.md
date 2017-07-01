# go-api

## Description

This project is based on https://github.com/ardanlabs/gotraining/blob/master/starter-kits/http.
It uses postgresql to persist data.

## Install dependencies

```sh
$ go get ./...
```

## Launch Docker containers (postgres)

You have to launch docker-compose in order to have all the required components up and running:

```sh
$ docker-compose up -d
```

## Build and run

```sh
$ go build ./cmd/apid
$ ./apid
```

## Database migrations

Migration use cli tool [migrate](https://github.com/mattes/migrate).

### Cli installation

```sh
$ go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
$ go build -tags 'postgres' -o /usr/local/bin/migrate github.com/mattes/migrate/cli
```

### Example

#### Create sql migrations files

```sh
$ migrate -database "postgres://go-api-postgres:go-api-postgres@localhost:5432/go-api-postgres?x-migrations-table=migrations" create -ext sql -dir migrations create_images
```

#### Exec migrations

```sh
$ migrate -source=file://migrations -database "postgres://go-api-postgres:go-api-postgres@localhost:5432/go-api-postgres?sslmode=disable" up
```

#### Rollback migrations

```sh
$ migrate -source=file://migrations -database "postgres://go-api-postgres:go-api-postgres@localhost:5432/go-api-postgres?sslmode=disable" down
```

## Swagger API documentation

You can access to the swagger API documentation at: http://[HOST][PORT]:3000/swagger/api-docs/

## Build Docker image from source

```sh
$ compile=CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo ./cmd/apid/
$ docker build -t go-api .
```

## Run Docker image

```sh
$ docker run --name go-api -e CONFIGOR_APPHOST=0.0.0.0 -e CONFIGOR_DATABASE_HOST=go-api-postgres --link go-api-postgres:go-api-postgres -p 3000:3000 -d go-api
```

## TODO

- [x] Migrations
- [x] Logger (Graylog)
- [ ] Make tests
- [x] Config (configor)
- [x] Swagger docs
- [x] Use sqlx instead of sql (structScan)
- [x] Health and readiness
- [ ] Prometheus metrics
- [x] Communicate with RabbitMQ
- [x] Dockerfile (and docker-compose)
- [ ] Jenkins integration
- [ ] govendor (production needs)

Thanks to contribute to this project. Each TODO must be done in a pull request.
