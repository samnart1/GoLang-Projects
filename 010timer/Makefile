# Go Simple Timer Makefile

# Variables
BINARY_NAME=timer
BUILD_DIR=build
MAIN_PACKAGE=./main.go

# Colors
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

.PHONY: help build run test clean install dev deps format lint

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-12s$(NC) %s\n", $$1, $$2}'

build: ## Build the application
	@echo "$(YELLOW)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

run: ## Run the application
	@echo "$(YELLOW)Running $(BINARY_NAME)...$(NC)"
	@go run $(MAIN_PACKAGE)

test: ## Run tests
	@echo "$(YELLOW)Running tests...$(NC)"
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning...$(NC)"
	@rm -rf $(BUILD_DIR)
	@go clean

install: build ## Install the binary to $GOPATH/bin
	@echo "$(YELLOW)Installing $(BINARY_NAME)...$(NC)"
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "$(GREEN)Installed to $(GOPATH)/bin/$(BINARY_NAME)$(NC)"

dev: ## Run in development mode with hot reload
	@echo "$(YELLOW)Starting development server...$(NC)"
	@go run $(MAIN_PACKAGE) --dev

deps: ## Download dependencies
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	@go mod download
	@go mod tidy

format: ## Format code
	@echo "$(YELLOW)Formatting code...$(NC)"
	@go fmt ./...

lint: ## Run linter
	@echo "$(YELLOW)Running linter...$(NC)"
	@golangci-lint run

# Build for multiple platforms
build-all: ## Build for multiple platforms
	@echo "$(YELLOW)Building for multiple platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	@echo "$(GREEN)Multi-platform build complete$(NC)"

# Quick commands
timer: build ## Build and run a 5-minute timer
	@$(BUILD_DIR)/$(BINARY_NAME) countdown -d 5m

stopwatch: build ## Build and run stopwatch
	@$(BUILD_DIR)/$(BINARY_NAME) stopwatch