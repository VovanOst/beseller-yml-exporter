.PHONY: help run build test clean fmt lint

help: ## Show this help
@echo "Available commands:"
@grep -E '^[a-zA-Z_-]+:.?## .' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf " \033[36m%-15s\033[0m %s\n", 1, $$2}'

run: ## Run the application
@echo "Running BeSeller YML Exporter..."
@go run cmd/exporter/main.go

build: ## Build the binary
@echo "Building..."
@mkdir -p bin
@go build -o bin/exporter cmd/exporter/main.go
@echo "Binary created: bin/exporter"

test: ## Run tests
@echo "Running tests..."
@go test -v -race -coverprofile=coverage.out ./...
@go tool cover -html=coverage.out -o coverage.html
@echo "Coverage report: coverage.html"

clean: ## Clean build artifacts
@echo "Cleaning..."
@rm -rf bin/
@rm -f coverage.out coverage.html
@echo "Clean complete"

fmt: ## Format code
@echo "Formatting code..."
@go fmt ./...
@echo "Format complete"

lint: ## Run linter
@echo "Running linter..."
@golangci-lint run ./...
@echo "Lint complete"

deps: ## Download dependencies
@echo "Downloading dependencies..."
@go mod download
@go mod tidy
@echo "Dependencies ready"