# Backend Docs Index

Topic-based documentation for the Go backend. Each file covers one concern.
The `docs` agent reads this index first to locate the right file before diving into source code.

| Topic | File | Source files covered |
|---|---|---|
| Startup lifecycle & dependency initialisation | [bootstrap.md](bootstrap.md) | `internal/bootstrap/bootstrap.go`, `internal/server/server.go`, `cmd/api/main.go` |
| Database connection & query patterns | [database.md](database.md) | `internal/infrastructure/database/postgres/db.go`, `internal/infrastructure/database/postgres/health_repository.go`, `internal/domain/health.go`, `internal/usecase/health_usecase.go` |
| Schema migrations (goose) | [migrations.md](migrations.md) | `cmd/migrate/main.go`, `internal/infrastructure/database/migrations/`, `Makefile` |
| HTTP routing & handler patterns | [routing.md](routing.md) | `internal/handler/handler.go`, `internal/handler/routes.go`, `internal/handler/hello_handler.go`, `internal/handler/health_handler.go`, `internal/server/server.go` |
| Integration testing with Testcontainers | [testing.md](testing.md) | `internal/infrastructure/database/postgres/health_repository_test.go`, `internal/transport/handlers/hello_handler_test.go` |
| Error handling conventions | [error-handling.md](error-handling.md) | `internal/repository/postgres/health_repository.go`, `internal/handler/health_handler.go`, `cmd/api/main.go` |
| Environment variables | [environment.md](environment.md) | `.env`, `internal/bootstrap/bootstrap.go`, `internal/repository/postgres/db.go` |
