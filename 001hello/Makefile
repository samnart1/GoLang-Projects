.PHONY: build clean test lint help install

APP_NAME := hello-cli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -ldflags "-X github.com/yourusername/hello-cli/pkg/version.Version=$(VERSION) \
                     -X github.com/yourusername/hello-cli/pkg/version.GitCommit=$(GIT_COMMIT) \
                     -X github.com/yourusername/hello-cli/pkg/version.BuildDate=$(BUILD_DATE)"

## build: Build the application
build:
	@echo "Building $(APP_NAME)..."
	go build $(LDFLAGS) -o bin/$(APP_NAME) .

## install: Install the application
install:
	@echo "Installing $(APP_NAME)..."
	go install $(LDFLAGS) .

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean

## test: Run tests
test:
	@echo "Running tests..."
	go test -v ./...

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

## help: Show this help message
help: Makefile
	@echo "Choose a command:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'