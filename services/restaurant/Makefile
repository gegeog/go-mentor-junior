include .env
export

export SERVICE_ROOT=$(shell pwd)

run:
	@go run ${SERVICE_ROOT}/cmd/server/main.go

vet:
	@go vet ./...

fmt:
	@go fmt ./...

build:
	@go build ${SERVICE_ROOT}/cmd/server/main.go
