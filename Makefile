
include .env
export

DB_URL_LOCAL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable
NETWORK_NAME=subscription-service_default

start:
	@go run cmd/subscription-service/main.go

up:
	@docker compose up -d --build

down:
	@docker compose down


gen-swagger:
	@docker compose run --rm swagger init --parseDependency --parseInternal -g cmd/subscription-service/main.go


migrate-up:
	@docker compose run --rm migrate -path /migrations -database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@subscr-postgres:5432/$(POSTGRES_DB)?sslmode=disable" up

migrate-down:
	@docker compose run --rm migrate -path /migrations -database "$(DB_URL_LOCAL)" down 1

migrate-force:
	@docker compose run --rm migrate -path /migrations -database "$(DB_URL_LOCAL)" force $(version)

migrate-create:
	@docker compose run --rm migrate create -ext sql -dir /migrations -seq $(name)


env-up:
	@docker compose up -d subscr-postgres

env-down:
	@docker compose down subscr-postgres

env-cleanup:
	@read -p "Clear all environment files? The risk of data loss [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down subscr-postgres && \
		rm -rf out/pgdata && \
		echo "The environment files have been cleared."; \
	else \
		echo "Environment cleanup has been canceled."; \
	fi
