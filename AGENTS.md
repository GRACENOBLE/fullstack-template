# AGENTS.md

Fullstack project — Go + Gin backend, Next.js 16 + React 19 web app, PostgreSQL, Android mobile (Kotlin + Compose).

Each layer has its own `AGENTS.md` and `docs/` folder. Read the file closest to what you are editing:
- `backend/AGENTS.md` + `backend/docs/` — Go API
- `web/AGENTS.md` + `web/docs/` — Next.js UI
- `mobile/AGENTS.md` + `mobile/docs/` — Android app

Claude Code users: see `CLAUDE.md` for the feature development workflow and subagent definitions.

---

## Setup

Prerequisites: Go 1.21+, Node.js 20+, pnpm, Docker + Docker Compose, Android Studio Meerkat (2024.3+), JDK 17+, Android SDK API 36.

```bash
# 1. Database
cd backend && make docker-run

# 2. Backend API (hot reload, separate terminal)
cd backend && make watch        # → http://localhost:8080

# 3. Web app (separate terminal)
cd web && pnpm install && pnpm dev   # → http://localhost:3000

# 4. Mobile — open mobile/ in Android Studio and run on emulator/device
cd mobile && ./gradlew assembleDebug      # build only
cd mobile && ./gradlew installDebug       # build and install on connected device
```

---

## Testing

```bash
# Backend (Docker must be running)
cd backend && make test         # unit + integration
cd backend && make itest        # integration only

# Web
cd web && pnpm lint
cd web && pnpm build

# Mobile
cd mobile && ./gradlew lint
cd mobile && ./gradlew test                  # unit tests (no device needed)
cd mobile && ./gradlew connectedAndroidTest  # instrumented tests (emulator/device required)
```

All backend DB tests use **Testcontainers** (real PostgreSQL). Never mock the database.

---

## Code style

### Go
- Business logic in `internal/` only — `cmd/` wires things together, nothing more
- Return errors up the stack — never swallow them
- Parameterized SQL only (`$1`, `$2`) — no string-concatenated queries
- `go vet ./...` must pass

### TypeScript / React
- Strict TypeScript — no `any`
- Next.js App Router only — no `pages/` directory
- Server Components by default — `"use client"` only for browser APIs or React hooks
- Tailwind CSS v4 only — no CSS modules, no inline `style={}`

---

### Kotlin / Jetpack Compose (mobile)
- Single Activity — no Fragments; all navigation via Compose
- No logic in `@Composable` functions — state lives in ViewModels
- Material3 only — `androidx.compose.material3`; never import M2 `material`
- Theme tokens — `MaterialTheme.colorScheme.*` and `MaterialTheme.typography.*`; never hardcode colors
- Version catalog — all versions in `mobile/gradle/libs.versions.toml`
- `modifier: Modifier = Modifier` as the last defaulted parameter on all public Composables

---

## Security

- No hardcoded secrets or credentials anywhere in the codebase
- Never commit `.env` files or `local.properties`
- All SQL must use parameterized queries
- CORS `AllowOrigins` must not be `["*"]` in non-local environments
- No direct user input rendered without sanitization
- Mobile API keys go in `local.properties` (gitignored), exposed via `BuildConfig` only

---

## PR instructions

- Branch from `main` — no direct pushes to `main`
- Run `make test` (backend), `pnpm lint && pnpm build` (web), and `./gradlew lint && ./gradlew test` (mobile) before opening a PR
- One logical change per PR
- PR title: concise description of what changed, not implementation details
