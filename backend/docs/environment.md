---
topic: environment
last_verified: 2026-06-25
sources:
  - .env
  - .env.example
  - internal/bootstrap/bootstrap.go
  - internal/infrastructure/database/postgres/db.go
  - pkg/firebase/admin.go
  - .golangci.yml
  - Makefile
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
| `ENV` | `bootstrap.go` | — | Environment name (`local`, `production`) |
| `BLUEPRINT_DB_HOST` | `bootstrap.go` | — | Postgres host (**required**) |
| `BLUEPRINT_DB_PORT` | `bootstrap.go` | — | Postgres port (**required**) |
| `BLUEPRINT_DB_DATABASE` | `bootstrap.go` | — | Database name (**required**) |
| `BLUEPRINT_DB_USERNAME` | `bootstrap.go` | — | Postgres username (**required**) |
| `BLUEPRINT_DB_PASSWORD` | `bootstrap.go` | — | Postgres password (**required**) |
| `BLUEPRINT_DB_SCHEMA` | `bootstrap.go` | `public` | Postgres search_path schema |
| `BLUEPRINT_DB_SSLMODE` | `bootstrap.go` | `disable` | Postgres SSL mode (`disable`, `require`, `verify-full`) |
| `RATE_LIMIT_RPS` | `bootstrap.go` | `0` (disabled) | Max requests per second per IP. Set to `0` or omit to disable rate limiting. |
| `RATE_LIMIT_BURST` | `bootstrap.go` | `int(RPS) * 5`, min 1 | Token-bucket burst capacity. Derived as `int(RPS)*5` when omitted; clamped to 1 so fractional RPS values never block all traffic. |
| `FIREBASE_PROJECT_ID` | `bootstrap.go`, `pkg/firebase/admin.go` | — | Firebase project ID. When omitted the Firebase Admin client is not initialised and `FirebaseAuth` middleware is skipped (auth disabled). |
| `FIREBASE_SERVICE_ACCOUNT_JSON` | `bootstrap.go`, `pkg/firebase/admin.go` | — | Raw JSON content of a Firebase service account key file. When omitted the SDK falls back to Application Default Credentials (ADC) — appropriate for GCP-hosted deployments. Only relevant when `FIREBASE_PROJECT_ID` is set. |
| `REDIS_URL` | `bootstrap.go` | — | Redis connection URL. When omitted or empty, cache/Redis initialization is skipped and the app runs without Redis. |
| `BLUEPRINT_WS_ALLOWED_ORIGIN` | `internal/transport/handlers/ws_handler.go` | — | Allowed origin for WebSocket CORS checks in staging/production. When omitted, WebSocket origin validation is skipped (local dev). |
| `SENTRY_DSN` | `bootstrap.go`, `internal/transport/handlers/routes.go` | — | Sentry error-tracking DSN. When omitted, the Sentry middleware is not registered. |
| `MAILJET_API_KEY` | `bootstrap.go` | — | Mailjet API key. When omitted (or empty), `App.EmailSender` is `nil` and no email is sent. |
| `MAILJET_SECRET_KEY` | `bootstrap.go` | — | Mailjet secret key. Must be provided alongside `MAILJET_API_KEY`. |
| `FROM_EMAIL` | `bootstrap.go` | — | Verified Mailjet sender address (e.g. `no-reply@example.com`). Required when `MAILJET_API_KEY` and `MAILJET_SECRET_KEY` are set; startup fails if omitted. |
| `FROM_NAME` | `bootstrap.go` | — | Sender display name (e.g. `MyApp`). Only read when both Mailjet credentials are set. |
| `R2_ACCOUNT_ID` | `bootstrap.go` | — | Cloudflare account ID. When omitted, `App.StorageService` is `nil` and storage routes are not registered. |
| `R2_ACCESS_KEY` | `bootstrap.go` | — | R2 API token access key. Required when `R2_ACCOUNT_ID` is set; startup fails if omitted. |
| `R2_SECRET_KEY` | `bootstrap.go` | — | R2 API token secret key. Required when `R2_ACCOUNT_ID` is set; startup fails if omitted. |
| `R2_BUCKET` | `bootstrap.go` | — | R2 bucket name. Required when `R2_ACCOUNT_ID` is set; startup fails if omitted. |
| `R2_PUBLIC_URL` | `bootstrap.go` | — | Public base URL for the R2 bucket (custom domain or `r2.dev` subdomain). Required when `R2_ACCOUNT_ID` is set; startup fails if omitted. |
| `IPAPI_KEY` | `bootstrap.go` | — | ipapi.co API key (optional). Free tier works without a key; supplying one enables higher rate limits. |
| `CORS_ALLOWED_ORIGINS` | `bootstrap.go` | `http://localhost:3000` | Comma-separated list of origins allowed by the CORS middleware. Parsed at startup — each entry is whitespace-trimmed. E.g. `https://app.example.com,https://staging.example.com`. |

Variables marked **required** are validated by `bootstrap.validateConfig` at startup — the process exits before attempting a DB connection if any are missing.

## `.env` file
Located at `backend/.env`. Never commit this file with real credentials.
The `.gitignore` in `backend/` excludes `.env` (verify before committing).

Docker Compose reads the same `.env` file to configure the Postgres container, so the values must be consistent between the app and Docker.

## Quality commands

| Command | What it runs |
|---|---|
| `make test` | `go test ./... -v` |
| `make itest` | `go test ./internal/infrastructure/... -v` (requires Docker) |
| `make lint` | `golangci-lint run ./...` |
| `make swagger` | regenerates `docs/swagger/` from swaggo annotations |

`make lint` uses the config in `backend/.golangci.yml`. Enabled linters: `errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`, `gofmt`, `goimports`, `misspell`, `revive`, `bodyclose`, `noctx`, `exhaustive`. The `revive` `exported` rule is disabled. Run `make lint` locally before pushing; CI also runs it.

## Adding a new environment variable
1. Add to `backend/.env` with a descriptive name.
2. Read it in `internal/bootstrap/bootstrap.go` inside `loadConfig()` and store it on `Config`.
3. If required, add a `requireNonEmpty` call in `validateConfig`.
4. Document it in this file.
5. Update `docker-compose.yml` if Docker also needs it.
