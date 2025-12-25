# Adding New Features

## Adding a New Domain Entity (e.g., Project, Tag)

### 1. Create domain package: `internal/{entity}/`

```
internal/{entity}/
├── {entity}.go    # Types, validation
├── errors.go      # Domain-specific errors
├── store.go       # Store interface
└── service.go     # Business logic
```

### 2. Define types in `{entity}.go`

```go
package project

type Project struct {
    ID        int64     `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"createdAt"`
}

type CreateRequest struct {
    Name string `json:"name"`
}

func (r CreateRequest) Validate() error {
    if r.Name == "" {
        return ErrValidation("name is required")
    }
    return nil
}
```

### 3. Define interface in `store.go` (consumer defines interface)

```go
package project

type Store interface {
    Create(ctx context.Context, p *Project) error
    GetByID(ctx context.Context, id int64) (*Project, error)
    List(ctx context.Context) ([]Project, error)
}
```

### 4. Implement store in `internal/sqlite/{entity}.go`

```go
package sqlite

type ProjectStore struct {
    db *DB
}

func NewProjectStore(db *DB) *ProjectStore {
    return &ProjectStore{db: db}
}

func (s *ProjectStore) Create(ctx context.Context, p *project.Project) error {
    // Implementation using parameterized queries
}
```

### 5. Add HTTP handlers in `internal/http/{entity}.go`

### 6. Wire in `cmd/server/main.go`

### 7. Add migration in `migrations/00X_{description}.sql`

## Naming Checklist

- Package name is singular, lowercase: `project` not `projects`
- Type names read naturally with package: `project.Project`
- No `I` prefix on interfaces: `Store` not `IStore`
- Error constructors: `ErrValidation()`, `ErrNotFound()`

## Don't

- Create separate `dto/`, `model/`, `entity/` packages
- Add unnecessary abstraction layers
- Skip validation at service boundary
- Forget to check `RowsAffected()` for updates/deletes
