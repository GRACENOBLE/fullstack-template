# ADR 0002 — PostgreSQL 16 + goose migrations

**Status:** Accepted  
**Date:** 2026-06-25

## Context

The template needs a primary relational store that:
- Supports ACID transactions and foreign key constraints
- Has a wide managed-hosting ecosystem (Supabase, AWS RDS, Neon, Railway)
- Can be run locally via Docker Compose with no cloud dependency
- Integrates cleanly with Go's `database/sql` interface

Schema evolution must be versioned, repeatable, and reviewable in pull requests.

## Decision

Use **PostgreSQL 16** (Docker image `postgres:16`) as the primary database.

Use **goose v3** (via `github.com/pressly/goose/v3`) for schema migrations, SQL-only (no Go migration files). Migration files live in `backend/internal/infrastructure/database/migrations/` and are applied with `make migrate-up`.

Use **pgx v5** (`github.com/jackc/pgx/v5`) as the driver via its `stdlib` adapter so repository code only depends on `database/sql` interfaces, not pgx types.

## Consequences

### Positive
- SQL migrations are version-controlled, human-readable, and reviewable in PRs.
- `goose` supports up/down migrations, named versions, and status queries — all exposed via the `Makefile`.
- pgx v5 via `stdlib` means repository interfaces stay portable; switching drivers requires no business logic changes.
- PostgreSQL's JSONB columns, full-text search, and window functions are available if needed.
- Testcontainers spins a real Postgres container in CI, so integration tests run against the same engine as production.

### Negative / trade-offs
- Docker must be running locally for `make docker-run` and integration tests.
- Migrations cannot be rolled back automatically after data has been written; rollbacks require a new forward migration.
- No ORM — raw SQL requires more boilerplate but gives full control over query plans and avoids N+1 surprises.
