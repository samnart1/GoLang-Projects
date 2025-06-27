.PHONY: build clean test lint fmt install dev run-examples help

# Variables
BINARY_NAME=json-parser
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=./main.go
GO_FILES=$(shell find . -name "*.go" -type f)

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

# Build the application
build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Binary built: $(BINARY_PATH)"

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@go clean

# Run tests
test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Lint code
lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

# Format code
fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w $(GO_FILES)

# Install dependencies
install: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Development mode
dev: build ## Build and run development examples
	@echo "Running development examples..."
	@$(BINARY_PATH) parse examples/simple.json
	@$(BINARY_PATH) parse examples/complex.json --format pretty

# Run example commands
run-examples: build ## Run example commands
	@echo "=== Parsing simple JSON ==="
	@$(BINARY_PATH) parse examples/simple.json
	@echo ""
	@echo "=== Parsing with pretty format ==="
	@$(BINARY_PATH) parse examples/complex.json --format pretty
	@echo ""
	@echo "=== Validating JSON ==="
	@$(BINARY_PATH) validate examples/simple.json
	@echo ""
	@echo "=== Formatting JSON ==="
	@$(BINARY_PATH) format examples/complex.json --indent 2

# Install binary to system
install-binary: build ## Install binary to system PATH
	@echo "Installing binary to /usr/local/bin..."
	@sudo cp $(BINARY_PATH) /usr/local/bin/$(BINARY_NAME)
	@echo "Binary installed: /usr/local/bin/$(BINARY_NAME)"