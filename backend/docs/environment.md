---
topic: environment
last_verified: 2026-06-14
sources:
  - .env
  - internal/bootstrap/bootstrap.go
  - internal/repository/postgres/db.go
---

# Environment Variables

## Loading mechanism
`godotenv` is loaded automatically via a blank import in `internal/bootstrap/bootstrap.go`:
```go
_ "github.com/joho/godotenv/autoload"
```
This runs on package init before any env var is read — no explicit `godotenv.Load()` call needed. Because `bootstrap` is the first package imported in `main`, `.env` is loaded before config validation runs.

## Variables reference

| Variable | Used in | Default | Description |
|---|---|---|---|
| `PORT` | `bootstrap.go` | `8080` | HTTP server listen port |
| `APP_ENV` | `bootstrap.go` | — | Environment name (`local`, `production`) |
| `BLUEPRINT_DB_HOST` | `bootstrap.go` | — | Postgres host (**required**) |
| `BLUEPRINT_DB_PORT` | `bootstrap.go` | — | Postgres port (**required**) |
| `BLUEPRINT_DB_DATABASE` | `bootstrap.go` | — | Database name (**required**) |
| `BLUEPRINT_DB_USERNAME` | `bootstrap.go` | — | Postgres username (**required**) |
| `BLUEPRINT_DB_PASSWORD` | `bootstrap.go` | — | Postgres password (**required**) |
| `BLUEPRINT_DB_SCHEMA` | `bootstrap.go` | `public` | Postgres search_path schema |
| `BLUEPRINT_DB_SSLMODE` | `bootstrap.go` | `disable` | Postgres SSL mode (`disable`, `require`, `verify-full`) |

Variables marked **required** are validated by `bootstrap.validateConfig` at startup — the process exits before attempting a DB connection if any are missing.

## `.env` file
Located at `backend/.env`. Never commit this file with real credentials.
The `.gitignore` in `backend/` excludes `.env` (verify before committing).

Docker Compose reads the same `.env` file to configure the Postgres container, so the values must be consistent between the app and Docker.

## Adding a new environment variable
1. Add to `backend/.env` with a descriptive name.
2. Read it in `internal/bootstrap/bootstrap.go` inside `loadConfig()` and store it on `Config`.
3. If required, add a `requireNonEmpty` call in `validateConfig`.
4. Document it in this file.
5. Update `docker-compose.yml` if Docker also needs it.
