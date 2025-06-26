# Go Web Scraper

A powerful and flexible web scraper built with Go that can extract titles, descriptions, links, images, and other metadata from web pages.

## Features

- **Single URL Scraping**: Extract detailed information from individual web pages
- **Batch Processing**: Scrape multiple URLs concurrently from a file
- **HTTP API Server**: RESTful API for web scraping operations
- **Configurable Options**: Custom timeouts, user agents, headers, and more
- **Multiple Output Formats**: JSON and pretty-printed output
- **Rate Limiting**: Built-in rate limiting to be respectful to target servers
- **Concurrent Processing**: Configurable concurrency for batch operations

## Installation

### From Source

```bash
git clone https://github.com/your-username/go-web-scraper.git
cd go-web-scraper
make setup
make build
```

### Using Go Install

```bash
go install github.com/your-username/go-web-scraper@latest
```

## Usage

### Command Line Interface

#### Scrape a Single URL

```bash
# Basic scraping
./bin/webscraper scrape https://example.com

# JSON output
./bin/webscraper scrape https://example.com --format json

# Custom options
./bin/webscraper scrape https://example.com --timeout 60s --agent "Custom Bot"

# Exclude links and images
./bin/webscraper scrape https://example.com --no-links --no-images
```

#### Batch Scraping

```bash
# Scrape URLs from file
./bin/webscraper batch --input urls.txt

# Save results to file
./bin/webscraper batch --input urls.txt --output results.json

# Custom concurrency and rate limiting
./bin/webscraper batch --input urls.txt --concurrent 5 --rate-limit 200ms
```

#### HTTP Server

```bash
# Start server on default port (8080)
./bin/webscraper server

# Custom port and host
./bin/webscraper server --port 3000 --host 0.0.0.0
```

### HTTP API

#### Health Check

```bash
curl http://localhost:8080/health
```

#### Scrape Single URL

```bash
curl -X POST http://localhost:8080/scrape \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "include_links": true,
    "include_images": true,
    "timeout": "30s"
  }'
```

#### Batch Scraping

```bash
curl -X POST http://localhost:8080/scrape/batch \
  -H "Content-Type: application/json" \
  -d '{
    "urls": ["https://example.com", "https://google.com"],
    "include_links": true,
    "include_images": true,
    "concurrent": 3
  }'
```

## Configuration

The scraper can be configured using environment variables:

```bash
# HTTP Client settings
export SCRAPER_TIMEOUT=30s
export SCRAPER_USER_AGENT="Custom-Bot/1.0"

# Server settings
export SERVER_PORT=8080
export SERVER_HOST=localhost

# Scraping defaults
export MAX_CONCURRENT=10
export RETRY_ATTEMPTS=3
export RATE_LIMIT=100ms
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SCRAPER_TIMEOUT` | HTTP request timeout | `30s` |
| `SCRAPER_USER_AGENT` | User agent string | `Go-Web-Scraper/1.0` |
| `SERVER_PORT` | HTTP server port | `8080` |
| `SERVER_HOST` | HTTP server host | `localhost` |
| `MAX_CONCURRENT` | Max concurrent workers | `10` |
| `RETRY_ATTEMPTS` | Number of retry attempts | `3` |
| `RATE_LIMIT` | Delay between requests | `100ms` |

## Input File Format

For batch processing, create a text file with one URL per line:

```
# urls.txt
https://example.com
https://github.com
https://golang.org
# Comments are ignored
https://stackoverflow.com
```

## Output Format

### Single Scrape Result

```json
{
  "url": "https://example.com",
  "title": "Example Domain",
  "description": "This domain is for use in illustrative examples",
  "keywords": ["example", "domain"],
  "links": [
    {
      "url": "https://www.iana.org/domains/example",
      "text": "More information...",
      "rel": ""
    }
  ],
  "images": [
    {
      "url": "https://example.com/image.png",
      "alt": "Example image",
      "src": "/image.png"
    }
  ],
  "headers": {
    "h1": "Example Domain"
  },
  "status_code": 200,
  "success": true,
  "scraped_at": "2024-01-15T10:30:00Z",
  "duration": "500ms"
}
```

### Batch Results

```json
{
  "results": [...],
  "total": 10,
  "success": 8,