---
topic: testing
last_verified: 2026-06-14
sources:
  - internal/database/database_test.go
---

# Testing

## Philosophy
**No mocks for the database.** All DB tests run against a real PostgreSQL instance spun up by Testcontainers. This catches schema mismatches, query errors, and type coercion issues that mocks hide.

## Testcontainers setup
`mustStartPostgresContainer()` starts a `postgres:latest` container and sets the package-level connection variables (`host`, `port`, `database`, `username`, `password`) so `New()` connects to the test container.

```go
func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
    dbContainer, err := postgres.Run(
        context.Background(),
        "postgres:latest",
        postgres.WithDatabase("database"),
        postgres.WithUsername("user"),
        postgres.WithPassword("password"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(5*time.Second)),
    )
    // sets host, port, database, username, password package vars
    // returns dbContainer.Terminate
}
```

## TestMain pattern
All integration test files use `TestMain` to start/stop the container once per test run:

```go
func TestMain(m *testing.M) {
    teardown, err := mustStartPostgresContainer()
    if err != nil {
        log.Fatalf("could not start postgres container: %v", err)
    }
    m.Run()
    if teardown != nil && teardown(context.Background()) != nil {
        log.Fatalf("could not teardown postgres container: %v", err)
    }
}
```

## Package placement
Tests live in the **same package** as the code under test.
- DB tests: `package database` in `internal/database/`
- Server/handler tests (when added): `package server` in `internal/server/`

## Running tests
```bash
make test    # unit + integration (requires Docker)
make itest   # integration only
go test ./internal/database -v -run TestHealth  # single test
```

## Adding a new integration test
1. Add a `TestXxx(t *testing.T)` function in `database_test.go` (or a new `_test.go` file in the same package).
2. Call `New()` to get the service — it will connect to the Testcontainers instance.
3. Use table-driven tests for multiple cases:
```go
func TestGetUser(t *testing.T) {
    srv := New()
    tests := []struct {
        name    string
        id      int64
        wantErr bool
    }{
        {"valid user", 1, false},
        {"missing user", 999, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := srv.GetUser(context.Background(), tt.id)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Requirements
- Docker must be running for any integration test.
- `go test` runs all tests including integration. Use build tags if you need to separate them in the future.
