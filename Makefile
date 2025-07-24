BINARY_NAME=usdt-app

BUILD_FLAGS=

DB_URL ?= postgres://postgres:pass@localhost:5432/dbname?sslmode=disable
API_URL ?= https://grinex.io/api/v2/depth
PORT ?= 50051
GO := $(shell which go)

MAIN_PATH=cmd/app/main.go

MIGRATIONS_DIR=migrations

.PHONY: build test docker-build run lint

build:
	$(GO) mod download
	$(GO) build -o $(BINARY_NAME) $(MAIN_PATH)

test:
	$(GO) test -cover ./... -v

docker-build:
	docker build -t usdt-app .

run:
	@DB_URL="$(DB_URL)" API_URL="$(API_URL)" PORT="$(PORT)" $(GO) run $(MAIN_PATH)

lint:
	golangci-lint run

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

proto:
	protoc --go_out=. --go-grpc_out=. proto/*.proto
