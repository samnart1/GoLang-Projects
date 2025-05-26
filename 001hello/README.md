# Hello CLI

A production-ready "Hello, World!" CLI application built with Go, demonstrating best practices for Go project structure and tooling.

## Features

- Clean architecture with separation of concerns
- Cobra CLI framework for robust command-line interface
- Proper error handling and logging
- Build versioning and metadata
- Makefile for common tasks
- Production-ready project structure

## Installation

### From Source

```bash
git clone https://github.com/yourusername/hello-cli.git
cd hello-cli
make install
```

### Build Binary

```bash
make build
./bin/hello-cli
```

## Usage

```bash
# Basic usage
hello-cli

# Greet someone specific
hello-cli --name "Alice"
hello-cli -n "Bob"

# Verbose output
hello-cli --verbose
hello-cli -v -n "Charlie"

# Show version
hello-cli --version

# Show help
hello-cli --help
```

## Development

### Prerequisites

- Go 1.21 or later
- golangci-lint (for linting)

### Commands

```bash
# Build the application
make build

# Run tests
make test

# Run linter
make lint

# Clean build artifacts
make clean

# Show help
make help
```

### Project Structure

```
hello-cli/
├── cmd/              # CLI commands and configuration
├── internal/         # Private application code
│   └── app/         # Application logic
├── pkg/             # Public libraries
│   └── version/     # Version information
├── main.go          # Application entry point
├── go.mod           # Go module definition
├── Makefile         # Build automation
└── README.md        # Documentation
```

## Architecture

This project follows Go best practices:

- **`cmd/`**: Contains the CLI command definitions using Cobra
- **`internal/`**: Private application code that shouldn't be imported by other projects
- **`pkg/`**: Public libraries that could be imported by other projects
- **Clean separation**: Business logic is separated from CLI concerns
- **Dependency injection**: Configuration is passed to the application layer
- **Error handling**: Proper error propagation and handling
- **Logging**: Structured logging with configurable verbosity

## License

MIT License