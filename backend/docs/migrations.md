---
topic: migrations
last_verified: 2026-06-14
sources:
  - cmd/migrate/main.go
  - migrations/
  - Makefile
---

# Migrations

## Tool
goose v3 (`github.com/pressly/goose/v3`).
Entry point: `cmd/migrate/main.go` — a thin wrapper that reuses `postgres.NewPostgresDB` and reads the same `BLUEPRINT_DB_*` env vars as the server. No separate goose binary installation needed.

## File location
`backend/migrations/` — SQL files only. Naming: `YYYYMMDDHHMMSS_<slug>.sql`, created automatically by `make migrate-create`.

## Makefile targets
| Target | What it does |
|---|---|
| `make migrate-create name=<slug>` | Create a new timestamped SQL file in `migrations/` |
| `make migrate-status` | Show applied vs. pending migrations |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-up-one` | Apply only the next pending migration |
| `make migrate-down` | Roll back the last applied migration |
| `make migrate-down-to version=N` | Roll back to a specific version number |
| `make migrate-reset` | Roll back all migrations to version 0 |
| `make migrate-version` | Print the current schema version |

## SQL migration format
Each file must have exactly one `-- +goose Up` annotation. `-- +goose Down` is optional but should always be included.

```sql
-- +goose Up
CREATE TABLE users (
    id         BIGSERIAL    PRIMARY KEY,
    email      TEXT         NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE users;
```

Rules:
- Every statement must end with `;`
- DDL only — no application data mutations in migrations
- For multi-statement blocks (PL/pgSQL), wrap with `-- +goose StatementBegin` / `-- +goose StatementEnd`
- Avoid `-- +goose NO TRANSACTION` unless the statement genuinely cannot run in a transaction (e.g. `CREATE INDEX CONCURRENTLY`)

## Workflow for any new table
1. `make migrate-create name=add_<table>` — generates the timestamped file
2. Fill in `CREATE TABLE` (Up) and `DROP TABLE` (Down)
3. `make migrate-up` — applies to local DB
4. Build the repository layer against the new schema
5. `make itest` to verify integration tests pass

## Go migrations — not supported
Go migrations require the migration functions to be registered and compiled into the binary. `go run ./cmd/migrate` produces a fresh binary on each invocation with no registered functions. Use SQL migrations for all schema changes.

## Testcontainers and migrations
Repository integration tests do **not** run goose migrations. Testcontainers starts a blank Postgres instance; tests create their own schema via `testDB.Exec(...)` in `TestMain` or per-test setup. This keeps tests independent of migration history and fast.

## Goose tracking table
Goose creates `goose_db_version` in the schema set by `BLUEPRINT_DB_SCHEMA` (via `search_path` in the connection string). Never modify this table manually.

## Hard rules
- No DDL inside Go code — `CREATE TABLE` belongs in a migration file, not in a repository method
- No Go migration files — SQL only
- Never edit or delete an applied migration — add a new one to correct it
