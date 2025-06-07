# Makefile for Go HTTP Server

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Binary names
BINARY_NAME=go-http-server
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe

# Build info
VERSION?=1.0.0
BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S%z)
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Ldflags for build info
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

.PHONY: all build clean test coverage deps fmt vet run dev help

# Default target
all: test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) -v

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_WINDOWS) -v

# Build for all platforms
build-all: build-linux build-windows build

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Run the server
run: build
	@echo "Starting server..."
	./$(BINARY_NAME)

# Run in development mode with live reload (requires air)
dev:
	@echo "Starting development server with live reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Install it with: go install github.com/cosmtrek/air@latest"; \
		echo "Running without live reload..."; \
		$(GOCMD) run main.go; \
	fi

# Run server with specific port
run-port:
	@echo "Starting server on port $(PORT)..."
	PORT=$(PORT) ./$(BINARY_NAME)

# Create directories for project structure
setup:
	@echo "Setting up project structure..."
	mkdir -p cmd internal/{server,config,logger} pkg/{response,version}
	mkdir -p web/{static/{css,js},templates} docs scripts
	touch internal/server/{server.go,handlers.go,middleware.go,routes.go}
	touch internal/config/config.go internal/logger/logger.go
	touch pkg/response/json.go pkg/version/version.go
	touch web/templates/{index.html,about.html,error.html}
	touch web/static/css/style.css web/static/js/app.js
	touch .env.example README.md .gitignore

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOCMD) install github.com/cosmtrek/air@latest
	$(GOCMD) install golang.org/x/tools/cmd/goimports@latest

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --rm $(BINARY_NAME):$(VERSION)

# Generate documentation
docs:
	@echo "Generating documentation..."
	$(GOCMD) doc -all ./... > docs/api.md

# Check for security vulnerabilities
security:
	@echo "Checking for security vulnerabilities..."
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install it with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install it from https://golangci-lint.run/usage/install/"; \
	fi

# Check dependencies for updates
check-updates:
	@echo "Checking for dependency updates..."
	$(GOCMD) list -u -m all

# Initialize git repository and make first commit
git-init:
	@echo "Initializing git repository..."
	git init
	git add .
	git commit -m "Initial commit: Go HTTP Server project"

# Show help
help:
	@echo "Available commands:"
	@echo "  build         - Build the binary"
	@echo "  build-all     - Build for all platforms"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  coverage      - Run tests with coverage"
	@echo "  deps          - Download dependencies"
	@echo "  fmt           - Format code"
	@echo "  vet           - Run go vet"
	@echo "  run           - Build and run the server"
	@echo "  dev           - Run with live reload (requires air)"
	@echo "  setup         - Create project structure"
	@echo "  install-tools - Install development tools"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  lint          - Run linter"
	@echo "  security      - Check for security issues"
	@echo "  help          - Show this help message"