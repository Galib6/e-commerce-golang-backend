include .env
export


# Create a new migration file
# Usage:
#   make migration name=add_users_table    # explicit name
#   make migration                         # auto-generated timestamped name
migration:
	@# if name not provided, generate timestamped name
	@if [ -z "$(name)" ]; then \
		NAME=$$(date +%Y%m%d%H%M%S)_migration; \
	else \
		NAME=$(name); \
	fi; \
	echo "Creating migration: $$NAME"; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$NAME

# Docker fallback to create migration when migrate CLI is not installed
migration-docker:
	@# usage: make migration-docker name=optional_name
	@if [ -z "$(name)" ]; then \
		NAME=$$(date +%Y%m%d%H%M%S)_migration; \
	else \
		NAME=$(name); \
	fi; \
	echo "Creating migration with docker: $$NAME"; \
	docker run --rm -v "$$PWD"/$(MIGRATIONS_DIR):/migrations --workdir /migrations migrate/migrate create -ext sql -dir /migrations -seq $$NAME

# Run migration up
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose up

# Run down N migrations (default 1)
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose down $(N)


# Force database to a specific version (helpful for dirty state)
migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(V)

# Reset DB (DEV only: force version 0 then run all migrations)
migrate-reset:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force 0
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose up

migration-version:
	@echo "Checking current migration version..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" -verbose version


# Example usage:
# make migration name=initial_tables
# make migrate-up
# make migrate-down N=1
# make migrate-force V=0
# make migrate-reset


# ============================================
# Development Commands (like npm scripts)
# ============================================

# Run tests 
test:
	go test ./... -v

# Run tests with coverage
test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Generate swagger docs
swagger:
	swag init -g ./cmd/api/main.go -o ./docs

# Download dependencies (like: npm install)
deps:
	go mod download
	go mod tidy

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Vet code for errors
vet:
	go vet ./...

# Clean build artifacts
clean:
	rm -rf bin/ coverage.out coverage.html

# Start Postgres (Docker)
db-start:
	docker run --name grapi-postgres -e POSTGRES_PASSWORD=root -e POSTGRES_DB=golang_restapi_postgresql_test1 -p 5432:5432 -d postgres:15-alpine

# Stop and remove Postgres container
db-stop:
	docker stop grapi-postgres && docker rm grapi-postgres

# Docker compose up
docker-up:
	docker compose up --build -d

# Docker compose down
docker-down:
	docker compose down --volumes

# ============================================
# Quick Reference (run: make help)
# ============================================
help:
	@echo "Available commands:"
	@echo "  make run          - Run the server"
	@echo "  make build        - Build binary to bin/app"
	@echo "  make test         - Run all tests"
	@echo "  make test-cover   - Run tests with coverage report"
	@echo "  make swagger      - Generate Swagger docs"
	@echo "  make deps         - Download and tidy dependencies"
	@echo "  make lint         - Lint code (requires golangci-lint)"
	@echo "  make fmt          - Format code"
	@echo "  make vet          - Check for errors"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make db-start     - Start Postgres in Docker"
	@echo "  make db-stop      - Stop Postgres container"
	@echo "  make docker-up    - Start all services (docker compose)"
	@echo "  make docker-down  - Stop all services"
	@echo "  make migrate-up   - Run migrations"
	@echo "  make migrate-down N=1 - Rollback N migrations"
	@echo "  make migration name=xyz - Create new migration"

# Note:
# Makefiles require TAB for indentation. Spaces will break it silently.