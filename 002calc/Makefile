.PHONY: build clean test lint help install run-interactive run-calc

APP_NAME := calc
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS := -ldflags "-X github.com/samnart1/GoLang-Projects/tree/calc/002calc/pkg/version.Version=$(VERSION) \
                     -X github.com/samnart1/GoLang-Projects/tree/calc/002calc/pkg/version.GitCommit=$(GIT_COMMIT) \
                     -X github.com/samnart1/GoLang-Projects/tree/calc/002calc/pkg/version.BuildDate=$(BUILD_DATE)"

## build: Build the calculator application
build:
	@echo "Building $(APP_NAME)..."
	go build $(LDFLAGS) -o bin/$(APP_NAME) .

## install: Install the calculator application
install:
	@echo "Installing $(APP_NAME)..."
	go install $(LDFLAGS) .

## run-interactive: Run the calculator in interactive mode
run-interactive: build
	./bin/$(APP_NAME) interactive

## run-calc: Run a sample calculation
run-calc: build
	./bin/$(APP_NAME) calc "5 + 3"

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
	@echo "Calculator CLI Commands:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'