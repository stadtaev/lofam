# Lofam - Task Management Backend

A task management backend focused on engineering perfection, following idiomatic Go patterns.

## Project Structure

```
lofam/
├── cmd/server/main.go       # Application entry point, wiring
├── internal/
│   ├── http/
│   │   ├── server.go        # Router, middleware, response helpers
│   │   └── task.go          # Task HTTP handlers
│   ├── sqlite/
│   │   ├── db.go            # Database connection, migrations
│   │   └── task.go          # TaskStore implementation
│   └── task/
│       ├── task.go          # Domain types, validation
│       ├── errors.go        # Domain errors (ValidationError, NotFoundError)
│       ├── store.go         # Store interface (defined at consumer)
│       └── service.go       # Business logic
├── migrations/
│   └── 001_init.sql         # SQL schema reference
└── go.mod
```

## Architecture Principles

### Idiomatic Go Patterns

1. **Interface at consumer, not implementer**: `task.Store` interface is defined in `internal/task/store.go` (where it's used), not in `internal/sqlite/` (where it's implemented). This is canonical Go - accept interfaces, return structs.

2. **Package naming**: Short, lowercase, no underscores
   - `http` (aliased as `lofamhttp` in main to avoid stdlib conflict)
   - `sqlite`
   - `task`

3. **Type naming**: Package-qualified names read naturally
   - `task.Task`, `task.Status`, `task.CreateRequest`
   - NOT: `task.TaskEntity`, `task.TaskStatus`

4. **No enterprise patterns**: Avoid excessive layering
   - No `domain/`, `repository/`, `handler/` directories
   - No `ITaskRepository` interface naming
   - No DTO mapping between identical structures

### Dependency Flow

```
main.go
   ↓ creates
sqlite.DB → sqlite.TaskStore
   ↓ implements
task.Store (interface)
   ↓ injected into
task.Service
   ↓ injected into
http.Server
```

## Code Conventions

### Error Handling

- Domain errors defined in `task/errors.go` as typed structs
- Constructor functions: `ErrValidation(msg)`, `ErrNotFound(id)`
- HTTP layer uses `errors.As()` to map to status codes
- Unexpected errors logged and returned as 500

### Request/Response

- Request types: `CreateRequest`, `UpdateRequest` with `Validate()` methods
- Validation at service layer boundary
- JSON tags with `omitempty` for optional fields
- Pointers for optional update fields (`*string`, `*Status`)

### HTTP Handlers

- Methods on `*Server` struct
- Use `r.Context()` for context propagation
- Centralized error handling via `handleError()`
- Helper functions: `writeJSON()`, `writeError()`, `parseID()`

### Database

- Use `modernc.org/sqlite` (pure Go, no CGO)
- Single connection (`SetMaxOpenConns(1)`) for SQLite
- Parameterized queries (no SQL injection)
- Check `RowsAffected()` for update/delete operations

## API Endpoints

```
GET    /api/tasks      - List all tasks
POST   /api/tasks      - Create task
GET    /api/tasks/{id} - Get task by ID
PUT    /api/tasks/{id} - Update task
DELETE /api/tasks/{id} - Delete task
```

## Domain Model

### Task Status
- `todo` (default)
- `in_progress`
- `done`

### Task Priority
- `low`
- `medium` (default)
- `high`

## Configuration

Environment variables:
- `DB_PATH` - SQLite database file path (default: `lofam.db`)
- `PORT` - HTTP server port (default: `8080`)

## Development

```bash
# Run server
go run ./cmd/server

# Build
go build -o lofam ./cmd/server

# Test (when tests exist)
go test ./...
```

## Dependencies

Minimal, justified:
- `github.com/go-chi/chi/v5` - Lightweight router with middleware
- `modernc.org/sqlite` - Pure Go SQLite (no CGO required)
