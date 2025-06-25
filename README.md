# go-reader

A Go-based RSS feed reader and processing service that enhances RSS feeds with full-text content and provides a REST API for managing feeds and articles.

## Features

- **RSS Feed Management**: Add, retrieve, and delete RSS feeds
- **Full-Text Enhancement**: Integrates with FiveFilters Full-Text RSS service to extract complete article content
- **Content Processing**: Sanitizes and normalizes HTML content from articles
- **Automatic Updates**: Background fetcher updates all feeds every 5 minutes
- **REST API**: Simple HTTP API for managing feeds and retrieving articles
- **Database Storage**: PostgreSQL backend with proper data relationships

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.24.3+ (for local development)

### Running with Docker Compose

```bash
docker compose up
```

This starts:
- go-reader service on http://localhost:8080
- PostgreSQL database on port 5432
- FiveFilters Full-Text RSS service on http://localhost:8081

### Local Development

```bash
# Install dependencies
go mod download

# Run directly
go run main.go

# Or build binary
go build -o ./tmp/main .
./tmp/main
```

## API Usage

### Add RSS Feed

```bash
curl -X POST http://localhost:8080/api/rss \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com/rss"}'
```

### Get All RSS Feeds

```bash
curl http://localhost:8080/api/rss
```

### Get Specific RSS Feed

```bash
curl http://localhost:8080/api/rss?id=1
```

### Get Articles for RSS Feed

```bash
# Get latest 100 articles (default)
curl http://localhost:8080/api/article?id=1

# Get specific number of articles
curl http://localhost:8080/api/article?id=1&limit=50
```

### Delete RSS Feed

```bash
curl -X DELETE http://localhost:8080/api/rss?id=1
```

## Architecture

### Components

- **HTTP Server**: REST API endpoints for feed and article management
- **RSS Processor**: Fetches and processes RSS feeds with content enhancement
- **Content Sanitizer**: Cleans and normalizes HTML content using bluemonday
- **Database Layer**: PostgreSQL with connection pooling for data persistence
- **Background Fetcher**: Automatically updates all feeds every 5 minutes

### Database Schema

**RSS Table**:
- Stores RSS feed metadata (URL, title, description)
- Includes FiveFilters enhanced URL
- Unique constraint on original URL

**Article Table**:
- Stores individual articles from RSS feeds
- Foreign key relationship to RSS feeds
- Unique constraint on (RSS ID, article link)
- Cascade delete when RSS feed is removed

### External Services

- **FiveFilters Full-Text RSS**: Enhances RSS feeds by extracting full article content from linked pages
- **PostgreSQL**: Primary database for persistent storage

## Development

### Hot Reloading

The Docker Compose setup uses [Air](https://github.com/air-verse/air) for automatic code reloading during development. Changes to Go files will automatically rebuild and restart the service.

### Testing

Use the Bruno API collection in `bruno-go-reader/` directory for testing API endpoints. Bruno provides a simple way to test all available API operations.

### Database Connection

Default connection string: `postgres://postgres:postgres@postgres:5432`

Tables are automatically created on application startup if they don't exist.

## Configuration

### Environment Variables

The application uses hardcoded configuration for simplicity. Key settings:

- **Server Port**: 8080
- **Database**: PostgreSQL connection to `postgres:5432`
- **Fetch Interval**: 5 minutes
- **FiveFilters URL**: `http://fullfeedrss:80`

### Docker Services

- **go-reader**: Main application (port 8080)
- **postgres**: PostgreSQL database (port 5432)
- **fullfeedrss**: FiveFilters service (port 8081)

