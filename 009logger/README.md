# Go Simple Logger

A professional-grade logging utility written in Go that provides flexible message logging with timestamps, multiple output formats, and both CLI and HTTP server interfaces.

## Features

- **Multiple Log Levels**: DEBUG, INFO, WARN, ERROR, FATAL
- **Flexible Output**: File, console, or both
- **Timestamp Formatting**: Customizable timestamp formats
- **CLI Interface**: Easy-to-use command-line interface
- **HTTP Server**: Web interface for remote logging
- **Log Rotation**: Built-in log file rotation support
- **Real-time Monitoring**: Tail command for live log viewing
- **Configuration**: YAML, JSON, or environment variable configuration

## Installation

### From Source
```bash
git clone https://github.com/user/go-simple-logger.git
cd go-simple-logger
make build
make install
```

### Using Go
```bash
go install github.com/user/go-simple-logger@latest
```

## Usage

### Basic Logging
```bash
# Log a simple message
go-simple-logger log "Hello, World!"

# Log with specific level
go-simple-logger log --level error "Something went wrong"

# Log to specific file
go-simple-logger log --file /var/log/app.log "Application started"
```

### Server Mode
```bash
# Start HTTP server
go-simple-logger server --port 8080

# Send logs via HTTP
curl -X POST http://localhost:8080/log \
  -H "Content-Type: application/json" \
  -d '{"message": "API log message", "level": "info"}'
```

### Real-time Monitoring
```bash
# Tail log file
go-simple-logger tail /var/log/app.log

# Follow with filtering
go-simple-logger tail --level error /var/log/app.log
```

## Configuration

Create a `.env` file or use environment variables:

```env
LOG_FILE=/var/log/simple-logger.log
LOG_LEVEL=info
LOG_FORMAT=json
MAX_SIZE=100
MAX_BACKUPS=3
MAX_AGE=7
```

## API Documentation

### POST /log
Log a message via HTTP.

**Request Body:**
```json
{
  "message": "Your log message",
  "level": "info",
  "source": "optional-source"
}
```

**Response:**
```json
{
  "success": true,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Development

```bash
# Install dependencies
make deps

# Run in development mode
make dev

# Run tests
make test

# Format and lint
make fmt
make lint
```

## License

MIT License - see LICENSE file for details.