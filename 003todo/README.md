# Todo CLI

A feature-rich command-line todo list manager written in Go.

## Features

- ✅ Add, edit, complete, and remove tasks
- 🔍 Search and filter tasks
- 📊 Task statistics and analytics
- 🎨 Colored terminal output
- 📱 Table and list view formats
- 🏷️ Task priorities and tags
- 📅 Due date tracking
- 💾 JSON file storage with backup
- 🔄 Data migration support

## Installation

### From Source

```bash
git clone https://github.com/yourusername/todo-cli.git
cd todo-cli
make build
```

### Using Go Install

```bash
go install github.com/yourusername/todo-cli@latest
```

## Usage

### Basic Commands

```bash
# Add a new task
todo add "Buy groceries"

# List all tasks
todo list

# Mark task as completed
todo done 1

# Remove a task
todo remove 1

# Edit a task
todo edit 1 "Buy groceries and cook dinner"

# Search tasks
todo search "groceries"

# Show statistics
todo stats
```

### Advanced Usage

```bash
# List only pending tasks
todo list --pending

# List tasks in table format
todo list --table

# Filter by priority
todo list --priority high

# Sort tasks by priority
todo list --sort priority
```

## Project Structure

```
todo-cli/
├── main.go              # Application entry point
├── cmd/                 # CLI command handlers
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   ├── task/            # Task domain logic
│   ├── storage/         # Data persistence
│   └── ui/              # User interface
├── pkg/                 # Public packages
└── testdata/            # Test data
```

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for build automation)

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Format code
make fmt
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
make test-coverage
```

## Configuration

The application stores data in `~/.todo-cli/`:

- `tasks.json` - Main task storage
- `backups/` - Automatic backups
- Configuration is managed automatically

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see LICENSE file for details.
```

## .gitignore
```
# Binaries
build/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out
coverage.html

# Go workspace file
go.work

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Todo CLI data (for development)
.todo-cli/