# Backend Docs Index

Topic-based documentation for the Go backend. Each file covers one concern.
The `docs` agent reads this index first to locate the right file before diving into source code.

| Topic | File | Source files covered |
|---|---|---|
| Startup lifecycle & dependency initialisation | [bootstrap.md](bootstrap.md) | `internal/bootstrap/bootstrap.go`, `internal/server/server.go`, `cmd/api/main.go` |
| Database connection & query patterns | [database.md](database.md) | `internal/infrastructure/database/postgres/db.go`, `internal/infrastructure/database/postgres/health_repository.go`, `internal/domain/health.go`, `internal/usecase/health_usecase.go` |
| Schema migrations (goose) | [migrations.md](migrations.md) | `cmd/migrate/main.go`, `internal/infrastructure/database/migrations/`, `Makefile` |
| HTTP routing & handler patterns | [routing.md](routing.md) | `internal/transport/handlers/handler.go`, `internal/transport/handlers/routes.go`, `internal/transport/handlers/hello_handler.go`, `internal/transport/handlers/health_handler.go`, `internal/transport/middleware/logger.go`, `internal/server/server.go` |
| Testing patterns (unit, handler, Redis, bootstrap) | [testing.md](testing.md) | `internal/infrastructure/database/postgres/health_repository_test.go`, `internal/transport/handlers/hello_handler_test.go`, `internal/transport/handlers/health_handler_test.go`, `internal/transport/middleware/logger_test.go`, `internal/usecase/health_usecase_test.go`, `internal/infrastructure/cache/redis/cache_test.go`, `internal/bootstrap/bootstrap_test.go` |
| Error handling conventions | [error-handling.md](error-handling.md) | `internal/infrastructure/database/postgres/health_repository.go`, `internal/transport/handlers/health_handler.go`, `cmd/api/main.go` |
| Environment variables | [environment.md](environment.md) | `.env`, `internal/bootstrap/bootstrap.go`, `internal/infrastructure/database/postgres/db.go` |
| Middleware (logger, rate limiter) | [middleware.md](middleware.md) | `internal/transport/middleware/logger.go`, `internal/transport/middleware/ratelimit.go`, `internal/transport/handlers/routes.go` |
| Firebase Auth (token verification, middleware, MeHandler) | [auth.md](auth.md) | `internal/usecase/auth_usecase.go`, `internal/transport/middleware/auth.go`, `internal/transport/handlers/auth_handler.go`, `pkg/firebase/admin.go`, `internal/bootstrap/bootstrap.go` |
| Observability (Sentry error tracking) | [observability.md](observability.md) | `internal/transport/middleware/sentry.go`, `internal/bootstrap/bootstrap.go`, `internal/transport/handlers/routes.go` |
