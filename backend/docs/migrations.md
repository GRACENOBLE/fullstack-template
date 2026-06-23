---
topic: migrations
last_verified: 2026-06-23
sources:
  - cmd/migrate/main.go
  - internal/infrastructure/database/migrations/
  - Makefile
---

# Migrations

## Tool
goose v3 (`github.com/pressly/goose/v3`).
Entry point: `cmd/migrate/main.go` — a thin wrapper that reuses `postgres.NewPostgresDB` and reads the same `BLUEPRINT_DB_*` env vars as the server. No separate goose binary installation needed.

## File location
`backend/internal/infrastructure/database/migrations/` — SQL files only. Naming: `YYYYMMDDHHMMSS_<slug>.sql`, created automatically by `make migrate-create`.

## Makefile targets
| Target | What it does |
|---|---|
| `make migrate-create name=<slug>` | Create a new timestamped SQL file in `internal/infrastructure/database/migrations/` |
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
- DDL only — no application data mutations in migrations, with one exception: `UPDATE` backfills are permitted as a step between `ADD COLUMN` (nullable) and `SET NOT NULL` (see pattern below)
- For multi-statement blocks (PL/pgSQL), wrap with `-- +goose StatementBegin` / `-- +goose StatementEnd`
- Avoid `-- +goose NO TRANSACTION` unless the statement genuinely cannot run in a transaction (e.g. `CREATE INDEX CONCURRENTLY`)

## Adding a NOT NULL column to an existing table
When adding a `NOT NULL` column to a table that may already have rows, follow the three-step pattern from `20260623063851_add_user_profile.sql`:

```sql
-- +goose Up
-- Step 1: add as nullable
ALTER TABLE users ADD COLUMN IF NOT EXISTS firebase_uid TEXT;

-- Step 2: backfill existing rows so SET NOT NULL will succeed
UPDATE users SET firebase_uid = 'legacy-' || id::text WHERE firebase_uid IS NULL;

-- Step 3: apply the constraint and unique index
ALTER TABLE users ALTER COLUMN firebase_uid SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT users_firebase_uid_key UNIQUE (firebase_uid);

-- +goose Down
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_firebase_uid_key;
ALTER TABLE users DROP COLUMN IF EXISTS firebase_uid;
```

This avoids the `ERROR: column contains null values` failure that occurs when you add a `NOT NULL` column with no default and the table is non-empty.

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
