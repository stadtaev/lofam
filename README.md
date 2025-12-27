# Lofam

A lightweight task management application with a Go backend and React frontend.

## Features

- Task management with status, priority, and due dates
- Self-hosted with SQLite storage (no external dependencies)
- RESTful API
- Modern React UI with TanStack Router

## Quick Start

### Using Docker Compose (Recommended)

```bash
docker compose up --build
```

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

### Development with Auto-Rebuild

```bash
docker compose watch
```

Automatically rebuilds when source files change.

### Local Development

**Backend:**
```bash
cd backend
go run ./cmd/server
```

**Frontend:**
```bash
cd frontend
bun install
bun dev
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/tasks | List all tasks |
| POST | /api/tasks | Create task |
| GET | /api/tasks/{id} | Get task by ID |
| PUT | /api/tasks/{id} | Update task |
| DELETE | /api/tasks/{id} | Delete task |

### Task Model

```json
{
  "id": 1,
  "title": "Task title",
  "description": "Optional description",
  "status": "todo | in_progress | done",
  "priority": "low | medium | high",
  "dueDate": "2025-12-25T00:00:00Z",
  "createdAt": "2025-12-27T10:00:00Z"
}
```

## Project Structure

```
lofam/
├── backend/              # Go API server
│   ├── cmd/server/       # Entry point
│   ├── internal/
│   │   ├── http/         # HTTP handlers
│   │   ├── sqlite/       # Database layer
│   │   └── task/         # Domain logic
│   └── Dockerfile
├── frontend/             # React SPA
│   ├── src/
│   │   ├── api/          # API client
│   │   ├── routes/       # TanStack Router pages
│   │   └── types/        # TypeScript types
│   ├── Dockerfile
│   └── nginx.conf
└── docker-compose.yml
```

## Tech Stack

**Backend:**
- Go 1.21+
- Chi router
- SQLite (modernc.org/sqlite)

**Frontend:**
- React 19
- TanStack Router
- Vite 7
- Pico CSS

## License

MIT
