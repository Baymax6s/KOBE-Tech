CMD_PATH := ./cmd/api
SWAG ?= swag
SWAG_VERSION ?= v1.16.4

.PHONY: run build test fmt tidy swagger swagger-check db-up db-down db-logs migrate-up migrate-down

run: swagger
	go run $(CMD_PATH)

build: swagger
	mkdir -p bin
	go build -o bin/api $(CMD_PATH)

swagger:
	@command -v $(SWAG) >/dev/null 2>&1 || (echo "swag is required. install with: go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)" && exit 1)
	$(SWAG) init -q -g ./cmd/api/main.go -d .,./internal --parseInternal -o ./swagger --ot json,yaml
	mv ./swagger/swagger.yaml ./swagger/openapi.yml

swagger-check: swagger
	git diff --exit-code -- swagger/openapi.yml swagger/swagger.json

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

migrate-up:
	docker compose run --rm --profile tools migrate up

migrate-down:
	docker compose run --rm --profile tools migrate down 1
