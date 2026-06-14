# Backend Docs Index

Topic-based documentation for the Go backend. Each file covers one concern.
The `docs` agent reads this index first to locate the right file before diving into source code.

| Topic | File | Source files covered |
|---|---|---|
| Database connection & query patterns | [database.md](database.md) | `internal/database/database.go` |
| HTTP routing & handler patterns | [routing.md](routing.md) | `internal/server/routes.go`, `internal/server/server.go` |
| Integration testing with Testcontainers | [testing.md](testing.md) | `internal/database/database_test.go` |
| Error handling conventions | [error-handling.md](error-handling.md) | `internal/database/database.go`, `cmd/api/main.go` |
| Environment variables | [environment.md](environment.md) | `.env`, `internal/database/database.go`, `internal/server/server.go` |
