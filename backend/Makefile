# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main.exe cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@docker compose up --build

# Shutdown DB container
docker-down:
	@docker compose down

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Integration tests (requires Docker)
itest:
	@echo "Running integration tests..."
	@go test ./internal/infrastructure/database/postgres/... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Output 'Watching...'; \
	} else { \
		Write-Output 'Installing air...'; \
		go install github.com/air-verse/air@latest; \
		air; \
		Write-Output 'Watching...'; \
	}"

# ---- Migrations (goose) ----

migrate-up:
	@go run ./cmd/migrate up

migrate-up-one:
	@go run ./cmd/migrate up-by-one

migrate-down:
	@go run ./cmd/migrate down

migrate-down-to:
	@go run ./cmd/migrate down-to $(version)

migrate-reset:
	@go run ./cmd/migrate down-to 0

migrate-status:
	@go run ./cmd/migrate status

migrate-version:
	@go run ./cmd/migrate version

migrate-create:
	@go run ./cmd/migrate create $(name)

.PHONY: all build run test clean watch docker-run docker-down itest \
        migrate-up migrate-up-one migrate-down migrate-down-to migrate-reset \
        migrate-status migrate-version migrate-create
