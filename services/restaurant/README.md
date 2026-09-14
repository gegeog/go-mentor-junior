# Restaurant service

Manages resraurants, menus, price and availbility.

## Prerequisites

* Go 1.27+
* Make

## Configuration

| Variable        | Default | Description              |
|-----------------|---------|--------------------------|
| `HTTP_ADDR`     |    —    | HTTP listen address (required), e.g. `:8080` |
| `LOGGER_LEVEL`  | `DEBUG` | Log level: `DEBUG`, `INFO`, `WARN`, `ERROR`  |

Configuration is read from `.env` file in the service root.

## Running locally 

```
cp .env.example .env # fill HTTP_ADDR
make run
```

## Commands
   Command                                                                          | Description
  ----------------------------------------------------------------------------------|----------------------------------------------------------------------------------
   make run                                                                         | Run the service
   make build                                                                       | Build the binary
   make fmt                                                                         | Format source code
   make vet                                                                         | Run static analysis
