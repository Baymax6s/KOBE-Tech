ifneq (,$(wildcard .env))
include .env
export
endif

MIGRATE_DATABASE_URL_DOCKER := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB)?sslmode=disable

.PHONY: setup dev test swagger db-up db-down migrate-down migrate-create

setup:
	@test -f .env || cp .env.example .env

dev: db-up
	air

swagger:
	swag init -q -g ./cmd/api/main.go -d .,./internal --parseInternal -o ./swagger --ot json,yaml
	mv ./swagger/swagger.yaml ./swagger/openapi.yml

db-up:
	@test -f .env || (echo "Run make setup first" >&2; exit 1)
	MIGRATE_DATABASE_URL="$(MIGRATE_DATABASE_URL_DOCKER)" docker compose up -d

db-down:
	docker compose down

migrate-down:
	@test -f .env || (echo "Run make setup first" >&2; exit 1)
	docker compose run --rm --entrypoint migrate migrate -path=/migrations -database="$(MIGRATE_DATABASE_URL_DOCKER)" down 1

migrate-create:
	@test -n "$(NAME)" || (echo "Usage: make migrate-create NAME=create_users" >&2; exit 1)
	docker compose run --rm --entrypoint migrate migrate create -ext sql -dir /migrations "$(NAME)"
