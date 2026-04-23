CMD_PATH := ./cmd/api

.PHONY: run build test fmt tidy db-up db-down db-logs migrate-up migrate-down migrate-create

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
	docker compose --profile tools run --rm migrate up

migrate-down:
	docker compose --profile tools run --rm migrate down 1

migrate-create:
	@test -n "$(NAME)" || (echo "Usage: make migrate-create NAME=create_users" && exit 1)
	docker compose --profile tools run --rm migrate create -ext sql -dir /migrations -seq $(NAME)
