.PHONY: build run clean test install dev

BINARY_NAME=guess-game
BUILD_DIR=bin

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

run: build
	@./$(BUILD_DIR)/$(BINARY_NAME) play

dev:
	@go run . play

clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

test:
	@echo "Running tests..."
	@go test ./...

install:
	@echo "Installing $(BINARY_NAME)..."
	@go install .
	@echo "Install complete"

tidy:
	@echo "Tidying dependencies..."
	@go mod tidy
	@echo "Dependencies updated"

# Add some useful development commands
play-easy:
	@go run . play -d easy

play-hard:
	@go run . play -d hard -t 30

play-timed:
	@go run . play -t 60

show-stats:
	@go run . stats

reset-data:
	@go run . reset
