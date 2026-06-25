# Architecture Decision Records

This directory captures significant architectural decisions made in the fullstack template. Each record documents the context, decision, and trade-offs at the time the choice was made.

## Format

Each ADR follows this structure:

- **Status** — Accepted / Deprecated / Superseded by [NNNN]
- **Date** — When the decision was recorded
- **Context** — The forces at play that drove a decision
- **Decision** — What was chosen and why
- **Consequences** — What becomes easier, what becomes harder

## Index

| # | Title | Status |
|---|---|---|
| [0001](0001-go-gin-backend.md) | Go 1.25 + Gin for the backend API | Accepted |
| [0002](0002-postgresql-goose.md) | PostgreSQL 16 + goose migrations | Accepted |
| [0003](0003-nextjs-app-router.md) | Next.js 16 App Router + React 19 | Accepted |
| [0004](0004-firebase-auth.md) | Firebase Authentication (cross-platform) | Accepted |
| [0005](0005-android-compose.md) | Android + Kotlin + Jetpack Compose + Material3 | Accepted |
| [0006](0006-redis-asynq.md) | Redis + Asynq for background jobs and event streaming | Accepted |
| [0007](0007-testcontainers.md) | Testcontainers for backend integration tests | Accepted |
| [0008](0008-cloudflare-r2.md) | Cloudflare R2 for object storage | Accepted |

## Adding a new ADR

1. Copy the next sequential number.
2. Create `NNNN-short-title.md` in this directory.
3. Add a row to the index table above.
4. Set **Status: Accepted** when the team aligns on the decision.
