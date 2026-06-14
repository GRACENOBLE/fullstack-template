---
topic: database
last_verified: 2026-06-14
sources:
  - internal/domain/health.go
  - internal/usecase/health_usecase.go
  - internal/repository/postgres/db.go
  - internal/repository/postgres/health_repository.go
---

# Database

## Driver
`database/sql` standard library with pgx v5 as the stdlib driver.
Import: `_ "github.com/jackc/pgx/v5/stdlib"` (blank import in `db.go`, registers the driver).
Do NOT use pgx's native API — use `database/sql` methods only.

## Connection
No singleton. `NewPostgresDB` is called once in `internal/server/server.go` and the returned `*sql.DB` is passed down the dependency chain. The caller is responsible for calling `db.Close()` on shutdown.

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

Env vars are loaded via `_ "github.com/joho/godotenv/autoload"` blank import in `internal/server/server.go` only.

## Architecture layers

```
internal/domain/health.go          — HealthStats type (map[string]string alias)
internal/usecase/health_usecase.go — HealthReader interface (repo contract), HealthUseCase interface, healthUseCase impl
internal/repository/postgres/      — HealthRepository: implements HealthReader against *sql.DB
```

The `HealthReader` interface is defined in the `usecase` package (Dependency Inversion — the use case owns the interface it depends on):

```go
// usecase/health_usecase.go
type HealthReader interface {
    Health(ctx context.Context) (domain.HealthStats, error)
}

type HealthUseCase interface {
    GetHealth(ctx context.Context) (domain.HealthStats, error)
}
```

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
3. Implement the repository in `internal/repository/postgres/` as a struct with a `New*` constructor.
4. Use `db.QueryContext` / `db.QueryRowContext` / `db.ExecContext`. Always pass `ctx`.
5. Always use parameterized queries — `$1`, `$2`, etc. Never string-concatenate SQL.
6. Return `(Result, error)` — never swallow errors or call `log.Fatal`.
7. Add integration test in `internal/repository/postgres/`.

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
