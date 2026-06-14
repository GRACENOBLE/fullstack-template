---
topic: testing
last_verified: 2026-06-15
sources:
  - internal/infrastructure/database/postgres/health_repository_test.go
  - internal/transport/handlers/hello_handler_test.go
---

# Testing

## Philosophy
**No mocks for the database.** All DB tests run against a real PostgreSQL instance spun up by Testcontainers. This catches schema mismatches, query errors, and type coercion issues that mocks hide.

## Testcontainers setup
`mustStartPostgresContainer()` starts a `postgres:latest` container, calls `NewPostgresDB(cfg)` with the container's mapped host/port, and assigns the result to the package-level `var testDB *sql.DB`.

```go
var testDB *sql.DB

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
    container, err := tcpostgres.Run(
        context.Background(),
        "postgres:latest",
        tcpostgres.WithDatabase("database"),
        tcpostgres.WithUsername("user"),
        tcpostgres.WithPassword("password"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(5*time.Second),
        ),
    )
    // resolves host and mapped port from container
    cfg := DBConfig{Host: dbHost, Port: dbPort.Port(), ...}
    db, err := NewPostgresDB(cfg)
    testDB = db
    return container.Terminate, nil
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
- Repository tests: `package postgres` in `internal/infrastructure/database/postgres/`
- Handler tests: `package handlers` in `internal/transport/handlers/`

## Handler unit tests
Handlers that have no DB dependency (e.g. `HelloWorldHandler`) use `httptest` without Testcontainers:

```go
func TestHelloWorldHandler(t *testing.T) {
    h := &Handler{}
    r := gin.New()
    r.GET("/", h.HelloWorldHandler)

    req, _ := http.NewRequest("GET", "/", nil)
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK { ... }
}
```

## Running tests
```bash
make test    # unit + integration (requires Docker)
make itest   # integration only — runs ./internal/infrastructure/database/postgres/...
go test ./internal/infrastructure/database/postgres/... -v -run TestHealth  # single test
```

## Adding a new integration test
1. Add a `TestXxx(t *testing.T)` function in a `_test.go` file under `internal/infrastructure/database/postgres/` (same package).
2. Construct the repository under test using `testDB`: e.g. `repo := NewHealthRepository(testDB)`.
3. Call repository methods directly and assert on the results.
4. Use table-driven tests for multiple cases:

```go
func TestGetUser(t *testing.T) {
    repo := NewUserRepository(testDB)
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
            _, err := repo.GetUser(context.Background(), tt.id)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Test ordering note
`internal/repository/postgres/health_repository_test.go` relies on test ordering: `TestNew` → `TestHealth` → `TestClose`. `TestClose` closes `testDB`; any test added after it that uses `testDB` will fail.

## Requirements
- Docker must be running for any integration test.
- `go test` runs all tests including integration. Use build tags if you need to separate them in the future.
