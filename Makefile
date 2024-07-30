# Makefile

# Load environment variables from .env file
include .env
export $(shell sed 's/=.*//' .env)

.PHONY: up create_migration migrate run stop

# Start PostgreSQL container
up:
	docker compose up -d postgres

# Create migration
create_migration:
	migrate create -ext sql -dir db/migrations -seq $(MIGRATION_NAME)

# Run migrations
migrate:
	migrate -path ./db/migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable up

# Start the application
run:
	docker compose up -d auth-service --build

# Stop all containers
stop:
	docker compose down
