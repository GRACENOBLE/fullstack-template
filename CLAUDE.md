# Project Guide

## Stack
| Layer | Tech |
|---|---|
| Frontend | Next.js 16, React 19, TypeScript 5, Tailwind CSS 4 |
| Backend | Go 1.25, Gin, `database/sql` + pgx v5 stdlib |
| Database | PostgreSQL 16 via Docker Compose |
| Dev tools | pnpm, Air (Go hot reload), Testcontainers |

## Run commands
```bash
cd backend && make docker-run   # start Postgres
cd backend && make watch        # backend dev with hot reload (Air) → :8080
cd frontend && pnpm dev         # frontend dev → :3000

cd backend && make test         # all tests
cd backend && make itest        # integration tests only (requires Docker)
cd frontend && pnpm lint && pnpm build
```

## Feature development workflow — always follow this
1. **Check docs first** — delegate to `docs` agent: find relevant topic docs, verify they match the current code
2. **Fix stale docs** — if docs diverge from code, update docs before implementing
3. **Implement** — delegate to `backend` or `frontend` agent, passing the relevant doc content as context
4. **Update docs** — delegate to `docs` agent: update `last_verified`, add new topics if introduced
5. **Quality gate** — run `/project:check` before declaring done

Use `/project:implement` to run this workflow end-to-end.

## Documentation locations
```
backend/docs/    # database, routing, testing, error-handling, environment
frontend/docs/   # routing, data-fetching, styling, components
```
Each doc file has `last_verified` and `sources` frontmatter. The `docs` agent maintains these.

## Available subagents — delegate to these
- **`backend`** — Go/Gin/PostgreSQL tasks
- **`frontend`** — Next.js/React/TypeScript tasks
- **`reviewer`** — pre-commit code review across both layers
- **`db-explorer`** — read-only DB schema and query analysis
- **`docs`** — documentation check, update, and creation

## Custom commands
- `/project:implement` — full documentation-first feature workflow
- `/project:check` — full quality gate (vet + lint + build + test)
- `/project:test` — run all tests with output
- `/project:new-route` — scaffold a new Go API route end-to-end

## Project layout
```
backend/
  cmd/api/main.go               # entry point — wiring only, no logic
  internal/
    server/server.go            # Server struct, NewServer()
    server/routes.go            # RegisterRoutes(), all handlers as *Server methods
    database/database.go        # Service interface + implementation, all queries
    database/database_test.go   # integration tests (Testcontainers)
  docker-compose.yml
  .env                          # never commit secrets
  Makefile
frontend/
  app/                          # Next.js App Router
  CLAUDE.md → AGENTS.md         # frontend-specific rules (read before writing Next.js)
```

## Go conventions
- Business logic in `internal/` only. `cmd/` just wires things together.
- New query → add to `Service` interface, implement on `service`, test in `database_test.go`.
- New route → register in `RegisterRoutes()`, handler as method on `*Server`.
- Return errors up the stack. Never `log.Fatal` or `os.Exit` inside `internal/`.
- Run `go vet ./...` before committing.

## TypeScript/React conventions
- App Router only. Default to Server Components.
- Add `"use client"` only when browser APIs or React hooks are required.
- Tailwind v4 for all styles — no CSS modules, no inline styles.
- No `any`. Use proper interfaces or `unknown`.
- Shared components → `components/`, utilities → `lib/`, types → `types/`.

## Testing — non-negotiable
- **Never mock the database.** Always use Testcontainers.
- Follow the `TestMain` + `mustStartPostgresContainer()` pattern in `database_test.go`.
- Tests live in the same package as the code (`package database`, `package server`).
- Docker must be running for integration tests.

## Hard rules (hooks enforce some of these)
- No secrets in committed files.
- No direct pushes to `main` — always open a PR.
- No `log.Fatal` / `os.Exit` inside `internal/`.
- No `"use client"` without a concrete browser requirement.
- No database mocks in tests.
- No `any` in TypeScript.
