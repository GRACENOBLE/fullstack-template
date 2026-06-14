---
topic: environment
last_verified: 2026-06-14
sources:
  - .env
  - internal/database/database.go
  - internal/server/server.go
---

# Environment Variables

## Loading mechanism
`godotenv` is loaded automatically via blank imports in two files:
- `internal/database/database.go`: `_ "github.com/joho/godotenv/autoload"`
- `internal/server/server.go`: `_ "github.com/joho/godotenv/autoload"`

This means `.env` in the working directory is loaded on package init — no explicit `godotenv.Load()` call needed.

## Variables reference

| Variable | Used in | Default | Description |
|---|---|---|---|
| `PORT` | `server.go` | `8080` | HTTP server listen port |
| `APP_ENV` | (available) | `local` | Environment name (`local`, `production`) |
| `BLUEPRINT_DB_HOST` | `database.go` | `localhost` | Postgres host |
| `BLUEPRINT_DB_PORT` | `database.go` | `5432` | Postgres port |
| `BLUEPRINT_DB_DATABASE` | `database.go` | `blueprint` | Database name |
| `BLUEPRINT_DB_USERNAME` | `database.go` | — | Postgres username |
| `BLUEPRINT_DB_PASSWORD` | `database.go` | — | Postgres password |
| `BLUEPRINT_DB_SCHEMA` | `database.go` | `public` | Postgres search_path schema |

## `.env` file
Located at `backend/.env`. Never commit this file with real credentials.
The `.gitignore` in `backend/` excludes `.env` (verify before committing).

Docker Compose reads the same `.env` file to configure the Postgres container, so the values must be consistent between the app and Docker.

## Adding a new environment variable
1. Add to `backend/.env` with a descriptive name.
2. Read with `os.Getenv("VAR_NAME")` at package level or inside the function that needs it.
3. Document it in this file.
4. Update `docker-compose.yml` if Docker also needs it.
