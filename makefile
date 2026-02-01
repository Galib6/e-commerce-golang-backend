include .env
export

# ============================================
# Migrations
# ============================================

GOPATH_BIN := $(shell go env GOPATH)/bin
ATLAS := $(GOPATH_BIN)/atlas

setup-tools:
	@echo "Installing Atlas..."
	@go install ariga.io/atlas/cmd/atlas@latest

# Generate migration by comparing models with database
# Usage: make migration name=add_user_phone
migration:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migration name=your_migration_name"; \
		exit 1; \
	fi
	@if [ ! -f $(ATLAS) ]; then $(MAKE) setup-tools; fi
	$(ATLAS) migrate diff $(name) --env local

# Apply all pending migrations (Using Go internal migrator)
migrate-up:
	@echo "🔄 Applying migrations..."
	@go run ./cmd/migrate

# Apply migrations using Atlas (Alternative)
migrate-up-atlas:
	@if [ ! -f $(ATLAS) ]; then $(MAKE) setup-tools; fi
	$(ATLAS) migrate apply --env local --url "$(DB_URL)"

# Show migration status
migrate-status:
	@if [ ! -f $(ATLAS) ]; then $(MAKE) setup-tools; fi
	$(ATLAS) migrate status --env local --url "$(DB_URL)"


# ============================================
# Development Commands
# ============================================

# Start dev server with hot-reload
dev:
	@echo "🚀 Starting dev server..."
	@air

# Start server (no hot-reload)
start:
	@go run ./cmd/api

# Build binary
build:
	@go build -o bin/app ./cmd/api
	@echo "✅ Built to bin/app"

# Run tests
test:
	@go test ./... -v

# Generate swagger docs
docs:
	@swag init -g ./cmd/api/main.go -o ./docs
	@echo "✅ Swagger docs generated"

# Install dependencies
install:
	@go mod download
	@go mod tidy

# Format code
fmt:
	@go fmt ./...

# Lint code
lint:
	@golangci-lint run


# ============================================
# Docker Commands
# ============================================

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down --volumes


# ============================================
# Help
# ============================================
help:
	@echo ""
	@echo "Available commands:"
	@echo ""
	@echo "  Development:"
	@echo "    make dev        - Start with hot-reload"
	@echo "    make start      - Start server"
	@echo "    make build      - Build binary"
	@echo ""
	@echo "  Database:"
	@echo "    make migration name=xxx  - Generate migration"
	@echo "    make migrate-up          - Run migrations"
	@echo "    make migrate-down N=1    - Rollback"
	@echo ""
	@echo "  Tools:"
	@echo "    make docs       - Generate Swagger"
	@echo "    make test       - Run tests"
	@echo "    make fmt        - Format code"
	@echo ""
