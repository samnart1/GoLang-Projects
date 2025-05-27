# Calculator CLI

A production-ready calculator CLI application built with Go, featuring both direct calculation and interactive modes.

## Features

- **Basic arithmetic operations**: +, -, *, /, ^, %
- **Two modes**: Direct calculation and interactive REPL
- **Decimal and negative number support**
- **Calculation history**
- **Comprehensive error handling**
- **Clean, modular architecture**

## Installation

```bash
git clone https://github.com/yourusername/calculator-cli.git
cd calculator-cli
make install
```

## Usage

### Direct Calculation Mode

```bash
# Basic operations
calc calc "5 + 3"        # Output: 5 + 3 = 8
calc calc "10 - 4"       # Output: 10 - 4 = 6
calc calc "7 * 6"        # Output: 7 * 6 = 42
calc calc "15 / 3"       # Output: 15 / 3 = 5

# Advanced operations
calc calc "2 ^ 8"        # Output: 2 ^ 8 = 256
calc calc "17 % 5"       # Output: 17 % 5 = 2

# Works with decimals and negatives
calc calc "3.14 * 2"     # Output: 3.14 * 2 = 6.28
calc calc "-5 + 10"      # Output: -5 + 10 = 5

# Verbose output
calc calc "10 / 3" --verbose
```

### Interactive Mode

```bash
calc interactive
# or
calc i

# Interactive session:
calc> 5 + 3
= 8
calc> 10 * 2.5
= 25
calc> help
calc> history
calc> exit
```

### Special Commands in Interactive Mode

- `help` - Show supported operations and examples
- `history` - Display calculation history
- `clear` - Clear calculation history
- `exit` or `quit` - Exit the calculator

## Architecture

```
calculator-cli/
├── cmd/                 # CLI commands
│   ├── root.go         # Root command setup
│   ├── calculate.go    # Direct calculation command
│   └── interactive.go  # Interactive mode command
├── internal/
│   ├── calculator/     # Core calculator logic
│   │   ├── calculator.go   # Main calculator struct
│   │   ├── operations.go   # Mathematical operations
│   │   └── parser.go       # Expression parsing
│   └── ui/             # User interface
│       └── interactive.go  # Interactive mode UI
├── pkg/
│   ├── version/        # Version information
│   └── errors/         # Custom error types
└── main.go            # Application entry point
```

## Development

```bash
# Build
make build

# Run tests
make test

# Quick test runs
make run-calc          # Test direct calculation
make run-interactive   # Test interactive mode

# Clean
make clean
```

## What's New in This Project

This calculator project introduces several new Go concepts:

1. **Regular expressions** for parsing mathematical expressions
2. **Custom error types** with context information  
3. **Interactive user input** with bufio.Scanner
4. **Mathematical operations** using the math package
5. **Data structures** like slices for history storage
6. **String manipulation** and formatting
7. **Multiple command structure** with Cobra subcommands

## Examples

```bash
# Simple calculations
calc calc "100 / 4"      # = 25
calc calc "3 ^ 4"        # = 81
calc calc "22 % 7"       # = 1

# Interactive session
calc interactive
calc> 15 + 25
= 40
calc> 2 ^ 10  
= 1024
calc> history
📊 Calculation History (2 entries):
1. 15 + 25 = 40
2. 2 ^ 10 = 1024
```