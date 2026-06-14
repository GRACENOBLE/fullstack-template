# AGENTS.md

Fullstack template — Go + Gin backend, Next.js 16 + React 19 frontend, PostgreSQL.

Each layer has its own `AGENTS.md` and `docs/` folder. Read the file closest to what you are editing:
- `backend/AGENTS.md` + `backend/docs/` — Go API
- `frontend/AGENTS.md` + `frontend/docs/` — Next.js UI

Claude Code users: see `CLAUDE.md` for the feature development workflow and subagent definitions.

---

## Setup

Prerequisites: Go 1.21+, Node.js 20+, pnpm, Docker + Docker Compose.

```bash
# 1. Database
cd backend && make docker-run

# 2. Backend API (hot reload, separate terminal)
cd backend && make watch        # → http://localhost:8080

# 3. Frontend (separate terminal)
cd frontend && pnpm install && pnpm dev   # → http://localhost:3000
```

---

## Testing

```bash
# Backend (Docker must be running)
cd backend && make test         # unit + integration
cd backend && make itest        # integration only

# Frontend
cd frontend && pnpm lint
cd frontend && pnpm build
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

## Security

- No hardcoded secrets or credentials anywhere in the codebase
- Never commit `.env` files
- All SQL must use parameterized queries
- CORS `AllowOrigins` must not be `["*"]` in non-local environments
- No direct user input rendered without sanitization

---

## PR instructions

- Branch from `main` — no direct pushes to `main`
- Run `make test` (backend) and `pnpm lint && pnpm build` (frontend) before opening a PR
- One logical change per PR
- PR title: concise description of what changed, not implementation details
