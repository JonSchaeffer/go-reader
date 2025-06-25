# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

go-reader is a Go-based RSS feed reader and processing service that:
- Fetches and manages RSS feeds
- Uses FiveFilters Full-Text RSS service to enhance feed content
- Processes article descriptions with HTML sanitization
- Stores feeds and articles in PostgreSQL database
- Provides REST API for managing RSS feeds and retrieving articles
- Runs background RSS fetcher every 5 minutes

## Development Commands

### Running the Application
- **Development with hot reload**: `docker compose up` (uses Air for hot reloading)
- **Direct Go run**: `go run main.go`
- **Build binary**: `go build -o ./tmp/main .`

### Database
- PostgreSQL runs on port 5432 (postgres/postgres credentials)
- Database connection: `postgres://postgres:postgres@postgres:5432`
- Tables are auto-created on startup (rss, article)

### Services
- Main app: http://localhost:8080
- FiveFilters RSS: http://localhost:8081
- PostgreSQL: localhost:5432

### Testing the API
Example Bruno API collection is available in `bruno-go-reader/` directory with test requests for:
- GET/POST/DELETE RSS feeds
- GET articles by RSS ID

## Architecture

### Core Modules

**`main.go`**: HTTP server setup, route handling, graceful shutdown
- Routes: `/api/rss` (GET/POST/DELETE), `/api/article` (GET)
- Background RSS fetcher starts automatically
- Database initialization and table creation

**`db/` package**: Database layer with connection pooling
- `connection.go`: pgx connection pool configuration (max 10, min 2 connections)
- `rss.go`: RSS feed CRUD operations, unique URL constraint
- `article.go`: Article storage with RSS foreign key, unique (rssID, link) constraint

**`rss/` package**: RSS processing and HTTP handlers
- `rss.go`: HTTP handlers, RSS feed fetching, FiveFilters integration
- `processor.go`: HTML content sanitization using bluemonday policy

### Key Data Flow
1. POST /api/rss → creates RSS entry → fetches via FiveFilters → saves articles
2. Background fetcher runs every 5 minutes → processes all RSS feeds → saves new articles
3. GET /api/article?id=X → returns articles for specific RSS feed

### External Dependencies
- **FiveFilters Full-Text RSS**: Enhances RSS feeds with full article content
- **PostgreSQL**: Primary data storage
- **bluemonday**: HTML sanitization for article content
- **Air**: Development hot-reloading tool

### Database Schema
- `rss` table: feeds with unique URL constraint, FiveFilters URL storage
- `article` table: articles with foreign key to RSS, unique (rssID, link) constraint
- Cascade deletes: removing RSS feed removes all its articles