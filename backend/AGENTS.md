# Backend AGENTS.md

Go backend: Gin + PostgreSQL (`database/sql` + pgx v5 stdlib) + Testcontainers.

---

## Commands

```bash
make docker-run    # start Postgres container
make docker-down   # stop Postgres container
make watch         # hot reload via Air → :8080
make run           # run once (no reload)
make build         # compile binary
make test          # all tests (requires Docker)
make itest         # integration tests only (requires Docker)
```

---

## Project structure

```
cmd/api/main.go                 # entry point — wiring and graceful shutdown only
internal/
  server/server.go              # Server struct, NewServer(), http.Server timeouts
  server/routes.go              # RegisterRoutes(), CORS, all handlers as *Server methods
  database/database.go          # Service interface, New() singleton, all query methods
  database/database_test.go     # Testcontainers integration tests
docker-compose.yml              # Postgres service (reads .env)
.env                            # env vars — never commit secrets
Makefile                        # all dev commands
```

---

## Detailed documentation

Read the relevant doc before implementing. These are kept in sync with the code — prefer them over general Go or Gin knowledge when they conflict.

| Topic | File |
|---|---|
| DB connection, Service interface, query patterns | [`docs/database.md`](docs/database.md) |
| Route registration, handler pattern, CORS | [`docs/routing.md`](docs/routing.md) |
| Testcontainers setup, TestMain, test patterns | [`docs/testing.md`](docs/testing.md) |
| Error handling, when log.Fatal is allowed | [`docs/error-handling.md`](docs/error-handling.md) |
| All environment variables and how they are loaded | [`docs/environment.md`](docs/environment.md) |

---

## Testing instructions

All database tests run against a **real PostgreSQL instance** via Testcontainers. Docker must be running.

```bash
make test     # unit + integration
make itest    # integration only
go test ./internal/database -v -run TestHealth   # single test
```

Adding a new test: follow the `TestMain` + `mustStartPostgresContainer()` pattern in `internal/database/database_test.go`.
Use table-driven tests for multiple input cases.

---

## Key conventions (short version — see docs for full detail)

- New route → register in `RegisterRoutes()`, handler as method on `*Server`
- New DB query → add to `Service` interface, implement on `*service`, add integration test
- Return errors up the stack — no `log.Fatal`/`os.Exit` inside `internal/` (see `docs/error-handling.md` for the two documented exceptions)
- Parameterized queries only — never `fmt.Sprintf` SQL
- Env vars loaded automatically via `godotenv/autoload` blank imports — no explicit call needed
