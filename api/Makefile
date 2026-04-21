CMD_PATH := ./cmd/api

.PHONY: run build test fmt tidy db-up db-down db-logs migrate-up migrate-down

run:
	go run $(CMD_PATH)

build:
	mkdir -p bin
	go build -o bin/api $(CMD_PATH)

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

migrate-up:
	docker compose run --rm --profile tools migrate up

migrate-down:
	docker compose run --rm --profile tools migrate down 1
