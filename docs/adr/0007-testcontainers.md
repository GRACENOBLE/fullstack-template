# ADR 0007 — Testcontainers for backend integration tests

**Status:** Accepted  
**Date:** 2026-06-25

## Context

Backend integration tests need to verify that repository code, SQL queries, and migrations work correctly against a real database engine. Options considered:
- Mock the database (in-memory fake or `sqlmock`)
- Use an in-process SQLite substitute
- Spin a real PostgreSQL container per test run (Testcontainers)

## Decision

Use **Testcontainers for Go** (`github.com/testcontainers/testcontainers-go`) to spin real PostgreSQL and Redis containers for integration tests.

Each integration test package follows the `TestMain` + `mustStartPostgresContainer()` pattern from `internal/infrastructure/database/postgres/health_repository_test.go`:

```go
func TestMain(m *testing.M) {
    db = mustStartPostgresContainer()
    os.Exit(m.Run())
}
```

Tests create their own schema via `testDB.Exec(...)` — they do **not** run goose migrations. This keeps tests independent of migration history and avoids state leakage between test runs.

**Database mocking is explicitly prohibited.** Handler tests may mock usecase interfaces; usecase tests may mock repository interfaces. Neither may mock the database itself.

## Consequences

### Positive
- Tests run against the exact same PostgreSQL version (16) as production — no mock/prod divergence.
- Container lifecycle is managed automatically; the Ryuk reaper cleans up orphaned containers.
- PostgreSQL-specific features (JSONB, `ON CONFLICT`, `RETURNING`) are tested correctly.
- CI (GitHub Actions, `ubuntu-latest`) runs Docker natively, so Testcontainers works without extra setup.

### Negative / trade-offs
- Docker must be running locally for integration tests. On Windows, rootless Docker mode is not supported by Testcontainers — contributors must use Docker Desktop.
- Integration tests are slower than unit tests (container startup ~1–2 s per package). The `make itest` target runs them separately from `make test` so the fast unit-test loop is not penalised.
- If Docker is unavailable (CI without Docker, or a locked-down corporate machine), integration tests cannot run.
