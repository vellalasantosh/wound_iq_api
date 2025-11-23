.PHONY: run test lint build

run:
	go run ./cmd/api

test:
	go test ./... -v

lint:
	@echo "Run `golangci-lint` or your preferred linter (not installed by default)."

build:
	go build -o bin/wound_iq_api ./cmd/api
