# Lofam - Task Management Application

A full-stack task management application with a Go backend and Next.js frontend.

## Project Structure

```
lofam/
├── backend/
│   ├── cmd/server/main.go       # Application entry point, wiring
│   ├── internal/
│   │   ├── http/
│   │   │   ├── server.go        # Router, middleware, CORS, response helpers
│   │   │   ├── task.go          # Task HTTP handlers
│   │   │   └── task_integration_test.go  # Integration tests
│   │   ├── sqlite/
│   │   │   ├── db.go            # Database connection, migrations
│   │   │   └── task.go          # TaskStore implementation
│   │   └── task/
│   │       ├── task.go          # Domain types, validation
│   │       ├── errors.go        # Domain errors (ValidationError, NotFoundError)
│   │       ├── store.go         # Store interface (defined at consumer)
│   │       └── service.go       # Business logic
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── app/
│   │   ├── layout.tsx           # Root layout with Inter font
│   │   ├── page.tsx             # Main calendar page (client component)
│   │   └── globals.css          # Tailwind CSS imports
│   ├── components/
│   │   ├── Calendar.tsx         # Month calendar with task indicators
│   │   ├── TaskList.tsx         # Tasks grouped by date with search
│   │   ├── TaskModal.tsx        # Create/edit/delete task modal
│   │   └── TodaySection.tsx     # Today's tasks + add button
│   ├── lib/
│   │   ├── api.ts               # Backend API client (fetch wrapper)
│   │   ├── types.ts             # TypeScript types (Task, CreateTaskRequest, etc.)
│   │   └── date-utils.ts        # Date helper functions
│   ├── Dockerfile               # Standalone Next.js build
│   ├── next.config.ts           # output: 'standalone' for Docker
│   └── package.json
├── infrastructure/
│   └── aws/
│       ├── provision.sh         # Idempotent AWS CLI provisioning
│       └── init-ssl.sh          # Let's Encrypt certificate setup
├── .github/workflows/
│   ├── deploy.yml               # CI/CD pipeline (auto on push)
│   └── infra.yml                # Infrastructure provisioning (manual)
├── docker-compose.yml           # Development (hot reload)
├── docker-compose.prod.yml      # Production (nginx + built images)
├── nginx.conf                   # Reverse proxy config
└── DEPLOYMENT.md                # AWS deployment guide
```

## Backend Architecture

### Idiomatic Go Patterns

1. **Interface at consumer, not implementer**: `task.Store` interface is defined in `internal/task/store.go` (where it's used), not in `internal/sqlite/` (where it's implemented).

2. **Package naming**: Short, lowercase, no underscores
   - `http` (aliased as `lofamhttp` in main to avoid stdlib conflict)
   - `sqlite`
   - `task`

3. **Type naming**: Package-qualified names read naturally
   - `task.Task`, `task.Status`, `task.CreateRequest`

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

### Code Conventions

**Error Handling:**
- Domain errors in `task/errors.go` as typed structs
- Constructor functions: `ErrValidation(msg)`, `ErrNotFound(id)`
- HTTP layer uses `errors.As()` to map to status codes

**HTTP:**
- CORS enabled via `go-chi/cors` middleware (allows localhost:3000)
- Methods on `*Server` struct
- Centralized error handling via `handleError()`

**Database:**
- `modernc.org/sqlite` (pure Go, no CGO)
- Parameterized queries (no SQL injection)

## Frontend Architecture

### Next.js 16 with App Router

- **Client components** for interactive UI (`'use client'`)
- **Tailwind CSS** for styling
- **Standalone output** for Docker deployment

### Component Structure

- `Calendar`: Month view with navigation, highlights dates with tasks
- `TaskList`: Searchable list of tasks grouped by due date
- `TaskModal`: Form for create/edit with status, priority, due date
- `TodaySection`: Quick view of today's tasks + add button

### API Client

- `lib/api.ts` wraps fetch calls to backend
- Uses `NEXT_PUBLIC_API_URL` env var (defaults to `http://localhost:8080`)
- Date format: RFC3339 (`2025-12-25T00:00:00Z`)

## API Endpoints

```
GET    /api/tasks      - List all tasks
POST   /api/tasks      - Create task
GET    /api/tasks/{id} - Get task by ID
PUT    /api/tasks/{id} - Update task
DELETE /api/tasks/{id} - Delete task
```

## Domain Model

### Task

```typescript
{
  id: number
  title: string
  description: string
  status: 'todo' | 'in_progress' | 'done'
  priority: 'low' | 'medium' | 'high'
  dueDate: string | null  // RFC3339 format
  createdAt: string       // RFC3339 format
}
```

## Configuration

### Backend Environment Variables
- `DB_PATH` - SQLite database file path (default: `lofam.db`)
- `PORT` - HTTP server port (default: `8080`)

### Frontend Environment Variables
- `NEXT_PUBLIC_API_URL` - Backend API URL (default: `http://localhost:8080`)

## Testing

### Integration Tests

Integration tests use in-memory SQLite and test the full HTTP stack.

```bash
cd backend
go test -tags=integration -v ./...
```

Test files use build tag `//go:build integration` to separate from unit tests.

**Test structure:**
- `internal/http/task_integration_test.go` - HTTP API integration tests
- Uses `httptest.Server` with real router
- In-memory SQLite (`:memory:`) as test double
- Table-driven tests for multiple scenarios

## Development

### Docker (Recommended)

```bash
# Build and run both services
docker compose up --build

# Auto-rebuild on file changes
docker compose watch
```

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

## Deployment

### Infrastructure (GitHub Actions)

Provision via GitHub Actions (manual trigger):

1. Create SSH key pair in AWS Console (EC2 → Key Pairs)
2. Add GitHub secrets:
   - `AWS_ACCESS_KEY_ID`
   - `AWS_SECRET_ACCESS_KEY`
   - `AWS_KEY_NAME` - key pair name
3. Run workflow: Actions → Infrastructure → Run workflow → provision

Creates (idempotent - checks before creating):
- EC2 t3.micro (free tier eligible)
- Security group (ports 22, 80, 443)
- Elastic IP

All resources tagged with `Project=lofam` for identification.

**Destroy**: Run workflow with "destroy" action.

### CI/CD (GitHub Actions)

On push to `main`:
1. Runs Go tests
2. SSHs to EC2
3. Pulls latest code
4. Runs `docker-compose -f docker-compose.prod.yml up --build -d`

**Required secrets:**
- `EC2_HOST` - Elastic IP address
- `EC2_SSH_KEY` - Private key content

### Manual Deployment

```bash
ssh -i ~/.ssh/your-key.pem ec2-user@<elastic-ip>
cd ~/app
git clone https://github.com/your/repo .  # First time
docker-compose -f docker-compose.prod.yml up --build -d
```

### Production Stack

```
nginx:80 → frontend:3000 (Next.js)
         → backend:8080  (Go API via /api/*)
```

## Dependencies

### Backend
- `github.com/go-chi/chi/v5` - Lightweight router with middleware
- `github.com/go-chi/cors` - CORS middleware
- `modernc.org/sqlite` - Pure Go SQLite (no CGO required)

### Frontend
- `next` 16.x - React framework
- `react` 19.x - UI library
- `tailwindcss` 4.x - Utility-first CSS
