include .env
export

export SERVICE_ROOT=$(shell pwd)

run:
	@go run ${SERVICE_ROOT}/cmd/server/main.go

build:
	@go build ${SERVICE_ROOT}/cmd/server/main.go

vet:
	@go vet ./...

fmt:
	@go fmt ./...

test:
	@go test ./...

test-race:
	@go test -race ./...

coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@rm -f coverage.out

generate:
	@go generate ${SERVICE_ROOT}/internal/api/handler/restaurant.go
	@go generate ${SERVICE_ROOT}/internal/service/restaurants.go

migrate-up:
	@docker compose run --rm migrator up

migrate-down:
	@docker compose run --rm migrator down

env-up:
	@docker compose up -d restaurants-postgres

env-down:
	@docker compose down restaurants-postgres

env-cleanup:
	@docker compose down -v restaurants-postgres
