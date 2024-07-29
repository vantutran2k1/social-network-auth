# Makefile

# Load environment variables from .env file
include .env
export $(shell sed 's/=.*//' .env)

.PHONY: up migrate run stop

# Start PostgreSQL container
up:
	docker-compose up -d postgres

# Run migrations
migrate:
	migrate -path ./db/migrations -database postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable up

# Start the application
run:
	docker-compose up -d auth-service --build

# Stop all containers
stop:
	docker-compose down
