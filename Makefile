.PHONY: help run build test clean migrate-up migrate-down lint fmt

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Run the application
	@echo "Starting application..."
	go run cmd/api/main.go

build: ## Build the application
	@echo "Building..."
	go build -o bin/my-notes-api cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

clean: ## Clean build files
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run ./...

fmt: ## Format code
	@echo "Formatting code..."
	gofmt -s -w .
	goimports -w .

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

migrate-create: ## Create new migration (usage: make migrate-create name=create_users_table)
	@if [ -z "$(name)" ]; then \
		echo "Error: name is required. Usage: make migrate-create name=create_users_table"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	@timestamp=$$(date +%Y%m%d%H%M%S); \
	touch migrations/$${timestamp}_$(name).up.sql; \
	touch migrations/$${timestamp}_$(name).down.sql; \
	echo "Created: migrations/$${timestamp}_$(name).up.sql"; \
	echo "Created: migrations/$${timestamp}_$(name).down.sql"

dev: ## Run with hot reload (requires air)
	@echo "Starting with hot reload..."
	air

install-tools: ## Install development tools
	@echo "Installing tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest
	@echo "Tools installed!"
