# Restaurant service

Manages restaurants, menus, price and availability.

## Prerequisites

* Go 1.27+
* Make

## Configuration

The process reads environment variables (`envconfig`). It does not load `.env` by itself.

`make run` / `make build` include `.env` and export it (see the Makefile). Copy `.env.example` first. Without Make, set the variables in the shell.

| Variable                | Default | Description                                         |
|-------------------------|---------|-----------------------------------------------------|
| `HTTP_ADDR`             | `:8080` | HTTP listen address                                 |
| `HTTP_SHUTDOWN_TIMEOUT` | `25s`   | Graceful shutdown deadline                          |
| `LOGGER_LEVEL`          | `DEBUG` | `DEBUG`, `INFO`, `WARN`, `ERROR` (case-insensitive) |

## Running locally

```bash
cp .env.example .env
make run
```

Or without Make:

```bash
HTTP_ADDR=:8080 go run ./cmd/server
```

## Testing

```bash
make test
```

With race detector:

```bash
make test-race
```

Coverage (prints per-function stats, cleans up after itself):

```bash
make coverage
```

## Mocks

Mocks are generated with [minimock](https://github.com/gojuno/minimock) pinned in `go.mod` as a tool dependency:

```bash
make generate
```

## Commands

| Command          | Description                          |
|------------------|--------------------------------------|
| `make run`       | Run the service                      |
| `make build`     | Build the binary                     |
| `make fmt`       | Format source code                   |
| `make vet`       | Run static analysis                  |
| `make test`      | Run tests                            |
| `make test-race` | Run tests with race detector         |
| `make coverage`  | Print statement coverage by function |
| `make generate`  | Regenerate mocks                     |
