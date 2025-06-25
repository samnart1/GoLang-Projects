# Makefile for Go File Reader

# Variables
APP_NAME := go-file-reader
BUILD_DIR := build
SRC_DIR := .
BIN := $(BUILD_DIR)/$(APP_NAME)
VERSION := dev
COMMIT := none
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Default target
.DEFAULT_GOAL := help

# Build the application
.PHONY: build
build: clean
	@echo "Building $(APP_NAME)..."
	@go build -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)" -o $(BIN) $(SRC_DIR)
	@echo "Build complete: $(BIN)"

# Run the application
.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	@$(BIN)

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean binaries and cache
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Clean complete."

# Display help
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  build     Build the application"
	@echo "  run       Build and run the application"
	@echo "  test      Run tests"
	@echo "  clean     Clean binaries and cache"
	@echo "  help      Display this help message"