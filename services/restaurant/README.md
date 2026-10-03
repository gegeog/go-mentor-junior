# Restaurant service

Manages resraurants, menus, price and availbility.

## Prerequisites

* Go 1.27+
* Make

## Configuration

The process reads environment variables (`envconfig`). It does not load `.env` by itself.

`make run` / `make build` include `.env` and export it (see the Makefile). Copy `.env.example` first. Without Make, set the variables in the shell.

| Variable                 | Default | Description |
|--------------------------|---------|-------------|
| `HTTP_ADDR`              | `:8080` | HTTP listen address |
| `HTTP_SHUTDOWN_TIMEOUT`  | `25s`   | Graceful shutdown deadline |
| `LOGGER_LEVEL`           | `DEBUG` | `DEBUG`, `INFO`, `WARN`, `ERROR` (case-insensitive) |

## Running locally 

```
cp .env.example .env
make run
```

Or without Make:

```
HTTP_ADDR=:8080 go run ./cmd/server
```

## Commands
   Command                                                                          | Description
  ----------------------------------------------------------------------------------|----------------------------------------------------------------------------------
   make run                                                                         | Run the service
   make build                                                                       | Build the binary
   make fmt                                                                         | Format source code
   make vet                                                                         | Run static analysis
