# Project Guide

## Stack
| Layer | Tech |
|---|---|
| Frontend | Next.js 16, React 19, TypeScript 5, Tailwind CSS 4 |
| Backend | Go 1.25, Gin, `database/sql` + pgx v5 stdlib |
| Database | PostgreSQL 16 via Docker Compose |
| Mobile | Android, Kotlin 2.2, Jetpack Compose BOM 2026.02, Material3 |
| Dev tools | pnpm, Air (Go hot reload), Testcontainers, Gradle 9.4 |

## Run commands
```bash
cd backend && make docker-run   # start Postgres
cd backend && make watch        # backend dev with hot reload (Air) → :8080
cd web && pnpm dev              # web dev → :3000

cd backend && make test         # all tests
cd backend && make itest        # integration tests only (requires Docker)
cd web && pnpm lint && pnpm build

cd mobile && ./gradlew assembleDebug          # build Android APK
cd mobile && ./gradlew lint && ./gradlew test # mobile quality gate
cd mobile && ./gradlew connectedAndroidTest   # instrumented tests (emulator/device required)
```

## Feature development workflow — always follow this
1. **Check docs first** — delegate to `docs` agent: find relevant topic docs, verify they match the current code
2. **Fix stale docs** — if docs diverge from code, update docs before implementing
3. **Implement** — delegate to `backend` or `web` agent, passing the relevant doc content as context
4. **Update docs** — delegate to `docs` agent: update `last_verified`, add new topics if introduced
5. **Quality gate** — run `/project:check` before declaring done

Use `/project:implement` to run this workflow end-to-end.

## Documentation locations
```
backend/docs/    # database, routing, testing, error-handling, environment
web/docs/        # routing, data-fetching, styling, components
mobile/docs/     # compose-conventions, architecture, testing
```
Each doc file has `last_verified` and `sources` frontmatter. The `docs` agent maintains these.

## Available subagents — delegate to these
- **`backend`** — Go/Gin/PostgreSQL tasks
- **`web`** — Next.js/React/TypeScript tasks
- **`mobile`** — Android/Kotlin/Jetpack Compose tasks
- **`reviewer`** — pre-commit code review across all layers
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
web/
  app/                          # Next.js App Router
  CLAUDE.md → AGENTS.md         # web-specific rules (read before writing Next.js)
mobile/
  app/src/main/java/com/company/template/
    MainActivity.kt             # single entry point, Compose root
    ui/theme/                   # Color, Theme, Type (Material3)
  gradle/libs.versions.toml     # version catalog — all versions declared here
  CLAUDE.md → AGENTS.md         # mobile-specific rules (read before writing Kotlin/Compose)
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

## Kotlin/Compose conventions
- Single Activity only. All navigation is Compose-based — no Fragments.
- No logic in `@Composable` functions — hoist state to ViewModel or the calling composable.
- Use Material3 (`androidx.compose.material3`) — not the older M2 `material` package.
- Theme tokens only — use `MaterialTheme.colorScheme.*` and `MaterialTheme.typography.*`; never hardcode colors in screens.
- All dependency versions declared in `mobile/gradle/libs.versions.toml`.
- Accept `modifier: Modifier = Modifier` as the last defaulted parameter in all public Composables.

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
