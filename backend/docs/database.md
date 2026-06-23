---
topic: database
last_verified: 2026-06-23
sources:
  - internal/domain/health.go
  - internal/domain/user.go
  - internal/domain/pagination.go
  - internal/usecase/health_usecase.go
  - internal/usecase/user.go
  - internal/infrastructure/database/postgres/db.go
  - internal/infrastructure/database/postgres/health_repository.go
  - internal/infrastructure/database/postgres/user_repository.go
---

# Database

## Driver
`database/sql` standard library with pgx v5 as the stdlib driver.
Import: `_ "github.com/jackc/pgx/v5/stdlib"` (blank import in `db.go`, registers the driver).
Do NOT use pgx's native API — use `database/sql` methods only.

## Connection
No singleton. `NewPostgresDB` is called once in `internal/bootstrap/bootstrap.go` and the returned `*sql.DB` is stored on `bootstrap.App`, then passed to repositories in `internal/server/server.go`. The caller is responsible for calling `db.Close()` on shutdown.

`DBConfig` holds all connection parameters:

```go
type DBConfig struct {
    Host     string
    Port     string
    Database string
    Username string
    Password string
    Schema   string
    SSLMode  string
}
```

`NewPostgresDB(cfg DBConfig) (*sql.DB, error)` builds the connection string and calls `sql.Open`. It returns an error instead of calling `log.Fatal`.

Connection string format:
```
postgres://username:password@host:port/database?sslmode=disable&search_path=schema
```

Env vars are loaded via `_ "github.com/joho/godotenv/autoload"` blank import in `internal/bootstrap/bootstrap.go`.

## Architecture layers

```
internal/domain/health.go                        — HealthStats type
internal/domain/user.go                          — User type
internal/domain/pagination.go                    — Page[T], CursorPage[T], PageRequest
internal/usecase/health_usecase.go               — HealthReader interface (repo contract), HealthUseCase interface + impl
internal/usecase/user.go                         — UserRepository interface
internal/infrastructure/database/postgres/       — HealthRepository, UserRepository: implement interfaces against *sql.DB
```

Repository interfaces are defined in the `usecase` package (Dependency Inversion — the use case owns the interface it depends on):

```go
// usecase/health_usecase.go
type HealthReader interface {
    Health(ctx context.Context) (domain.HealthStats, error)
}

type HealthUseCase interface {
    GetHealth(ctx context.Context) (domain.HealthStats, error)
}

// usecase/user.go
type UserRepository interface {
    Upsert(ctx context.Context, u *domain.User) (*domain.User, error)
    DeleteByFirebaseUID(ctx context.Context, firebaseUID string) error
}
```

## Domain types

```go
// internal/domain/user.go
type User struct {
    ID          int64     `json:"id"`
    FirebaseUID string    `json:"firebase_uid"`
    Name        string    `json:"name"`
    Email       string    `json:"email"`
    PhotoURL    string    `json:"photo_url"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// internal/domain/pagination.go
type Page[T any] struct {
    Items    []T  `json:"items"`
    Total    int  `json:"total"`
    Page     int  `json:"page"`
    PageSize int  `json:"page_size"`
    HasMore  bool `json:"has_more"`
}

type CursorPage[T any] struct {
    Items      []T    `json:"items"`
    NextCursor string `json:"next_cursor,omitempty"`
    HasMore    bool   `json:"has_more"`
}

type PageRequest struct {
    Page     int `form:"page"      binding:"min=0"`
    PageSize int `form:"page_size" binding:"min=0,max=100"`
}

func (p *PageRequest) Defaults()    // fills zero values with page=1, size=20
func (p PageRequest) Offset() int   // returns SQL OFFSET value
```

## Users table schema
Added via migration `20260623063851_add_user_profile.sql`:

| Column | Type | Constraints |
|---|---|---|
| `id` | `BIGSERIAL` | PRIMARY KEY |
| `firebase_uid` | `TEXT` | NOT NULL, UNIQUE |
| `name` | `TEXT` | — |
| `email` | `TEXT` | — |
| `photo_url` | `TEXT` | — |
| `created_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |
| `updated_at` | `TIMESTAMPTZ` | NOT NULL DEFAULT now() |

## UserRepository
`internal/infrastructure/database/postgres/user_repository.go` — constructor returns the interface, not the concrete type:

```go
func NewUserRepository(db *sql.DB) usecase.UserRepository
```

`Upsert` uses `INSERT ... ON CONFLICT (firebase_uid) DO UPDATE` and returns the full row via `RETURNING`.
`DeleteByFirebaseUID` deletes by `firebase_uid`; does not error when the row does not exist (DELETE is idempotent).

## Repository pattern
Each repository is a struct that holds `*sql.DB` and is constructed with a `New*` function.

```go
type HealthRepository struct {
    db *sql.DB
}

func NewHealthRepository(db *sql.DB) *HealthRepository {
    return &HealthRepository{db: db}
}
```

## Adding a new query — exact pattern
1. Define a domain type in `internal/domain/` if needed.
2. Add a repository interface in the relevant `usecase/` file (the use case owns the interface).
3. Implement the repository in `internal/infrastructure/database/postgres/` as a struct with a `New*` constructor.
4. Use `db.QueryContext` / `db.QueryRowContext` / `db.ExecContext`. Always pass `ctx`.
5. Always use parameterized queries — `$1`, `$2`, etc. Never string-concatenate SQL.
6. Return `(Result, error)` — never swallow errors or call `log.Fatal`.
7. Add integration test in `internal/infrastructure/database/postgres/`.

```go
// Repository method
func (r *UserRepository) GetUser(ctx context.Context, id int64) (*domain.User, error) {
    row := r.db.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id = $1", id)
    var u domain.User
    if err := row.Scan(&u.ID, &u.Name); err != nil {
        return nil, fmt.Errorf("postgres: GetUser: %w", err)
    }
    return &u, nil
}
```

## Health check
`HealthRepository.Health(ctx)` returns `(domain.HealthStats, error)`.
On ping failure: sets `stats["status"] = "down"` and returns a non-nil error.
On success: sets `stats["status"] = "up"` plus connection pool stats.
The HTTP handler returns 503 when this method returns an error (see routing.md).

## Connection pool stats
`Health()` exposes: `open_connections`, `in_use`, `idle`, `wait_count`, `wait_duration`, `max_idle_closed`, `max_lifetime_closed`.
Thresholds in `Health()`: warn at 40 open connections, 1000 wait events.
