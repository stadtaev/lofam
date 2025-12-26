# Lofam

A lightweight backend service for independent private data storage and organization.

## Features

- Self-hosted data management
- RESTful API
- SQLite storage (no external dependencies)
- Minimal footprint

## Quick Start

### Using Docker Compose (Recommended)

```bash
docker compose up -d
```

The API will be available at `http://localhost:8080`.

```bash
# View logs
docker compose logs -f

# Stop
docker compose down
```

### Local Development

```bash
go run ./cmd/server
```

## License

MIT
