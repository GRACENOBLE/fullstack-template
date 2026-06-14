---
topic: database
last_verified: 2026-06-14
sources:
  - internal/database/database.go
  - internal/database/database_test.go
---

# Database

## Driver
`database/sql` standard library with pgx v5 as the stdlib driver.
Import: `_ "github.com/jackc/pgx/v5/stdlib"` (blank import, registers the driver).
Do NOT use pgx's native API — use `database/sql` methods only.

## Connection
Singleton pattern via package-level `dbInstance *service` var.
`New()` returns early if `dbInstance != nil`.

Connection string format:
```
postgres://username:password@host:port/database?sslmode=disable&search_path=schema
```

Env vars loaded automatically via `_ "github.com/joho/godotenv/autoload"` blank import in `database.go`.

## Service interface
All DB operations are defined as methods on the `Service` interface.
The `service` struct is the private implementation.

```go
type Service interface {
    Health() map[string]string
    Close() error
    // Add new methods here
}

type service struct {
    db *sql.DB
}
```

## Adding a new query — exact pattern
1. Add method signature to `Service` interface.
2. Implement on `*service` using `s.db.QueryContext` / `s.db.QueryRowContext` / `s.db.ExecContext`.
3. Always use parameterized queries — `$1`, `$2`, etc. Never string-concatenate SQL.
4. Add integration test in `database_test.go`.

```go
// Interface
GetUser(ctx context.Context, id int64) (*User, error)

// Implementation
func (s *service) GetUser(ctx context.Context, id int64) (*User, error) {
    row := s.db.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id = $1", id)
    var u User
    if err := row.Scan(&u.ID, &u.Name); err != nil {
        return nil, err
    }
    return &u, nil
}
```

## Health check
`Health()` returns a `map[string]string` with `status`, `message`, and connection pool stats.
Called by the `/health` route. Calls `s.db.PingContext` with a 1-second timeout.
Note: `Health()` calls `log.Fatalf` on ping failure — this is intentional (terminates on unrecoverable DB loss).

## Connection pool stats
`Health()` exposes: `open_connections`, `in_use`, `idle`, `wait_count`, `wait_duration`, `max_idle_closed`, `max_lifetime_closed`.
Thresholds in `Health()`: warn at 40 open connections, 1000 wait events.
