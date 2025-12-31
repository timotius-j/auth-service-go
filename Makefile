APP_NAME=app

ifneq (,$(wildcard .env))
include .env
export
endif

MIGRATION_PATH=migrations
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: install dev run build test clean migrate-up migrate-down reset-db

install:
	go mod tidy

dev:
	air

run:
	go run ./cmd/server/main.go

build:
	go build -o bin/$(APP_NAME) ./cmd/server

test:
	go test ./... -v

clean:
	rm -rf tmp bin

migrate-up:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATION_PATH) -database "$(DB_URL)" down

reset-db:
	psql "$(DB_URL)" -f migrations/000_drop_tables.sql
