.PHONY: build clean run test install dev deps fmt lint

# Variables
BINARY_NAME=go-simple-logger
MAIN_PATH=./main.go
BUILD_DIR=./bin

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Run the application
run:
	@go run $(MAIN_PATH)

# Run with log command
run-log:
	@go run $(MAIN_PATH) log "Test message"

# Run with server command
run-server:
	@go run $(MAIN_PATH) server

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Development mode (with auto-rebuild)
dev:
	@echo "Starting development mode..."
	@which air > /dev/null || go install github.com/cosmtrek/air@latest
	@air

# Help
help:
	@echo "Available targets:"
	@echo "  build     - Build the application"
	@echo "  clean     - Clean build artifacts"
	@echo "  run       - Run the application"
	@echo "  run-log   - Run with log command"
	@echo "  run-server - Run with server command"
	@echo "  test      - Run tests"
	@echo "  deps      - Install dependencies"
	@echo "  fmt       - Format code"
	@echo "  lint      - Lint code"
	@echo "  install   - Install binary to /usr/local/bin"
	@echo "  dev       - Start development mode with auto-reload"