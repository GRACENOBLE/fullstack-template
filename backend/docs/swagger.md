---
topic: swagger
last_verified: 2026-06-15
sources:
  - cmd/api/main.go
  - internal/transport/handlers/routes.go
  - internal/transport/handlers/swagger_types.go
  - docs/swagger/docs.go
  - Makefile
---

# Swagger / OpenAPI Documentation

## Overview

Swagger UI is served at `GET /swagger/index.html` (proxied via `/swagger/*any`). The generated spec files live in `docs/swagger/` and are committed to version control.

**Tool:** [swaggo/swag](https://github.com/swaggo/swag) v1.16.x  
**UI middleware:** `github.com/swaggo/gin-swagger` + `github.com/swaggo/files`

## Regenerating docs

After adding or changing any handler annotation, run from `backend/`:

```bash
make swagger
```

This runs:
```bash
go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/api/main.go -o docs/swagger
```

Commit the updated `docs/swagger/docs.go`, `docs/swagger/swagger.json`, and `docs/swagger/swagger.yaml` alongside the handler change.

## Annotation locations

### API-level metadata — `cmd/api/main.go`

Place the block immediately above `func main()`:

```go
//	@title		Blueprint API
//	@version	1.0
//	@description	Fullstack template REST API.
//
//	@host		localhost:8080
//	@BasePath	/
func main() {
```

### Handler-level annotations

Place the block immediately above the handler method:

```go
//	@Summary	Short description
//	@Tags		tagname
//	@Produce	json
//	@Param		id	path	int	true	"Record ID"
//	@Success	200	{object}	ResponseType
//	@Failure	400	{object}	map[string]string
//	@Failure	500	{object}	map[string]string
//	@Router		/resource/{id} [get]
func (h *Handler) MyHandler(c *gin.Context) {
```

## Referencing domain types

Swaggo resolves types from the handler file's own imports. The `handlers` package does **not** directly import `backend/internal/domain`, so you cannot write `{object} domain.HealthStats` in an annotation.

**Fix:** add a type alias to `internal/transport/handlers/swagger_types.go`:

```go
// swagger_types.go
package handlers

import "backend/internal/domain"

type HealthStats = domain.HealthStats
// Add more aliases here as new domain types are introduced
```

Then reference the alias in the annotation:

```go
//	@Success	200	{object}	HealthStats
```

## Common `@Param` sources

| Source | swag keyword |
|--------|-------------|
| URL path | `path` |
| Query string | `query` |
| Request body | `body` |
| Header | `header` |

Example with a body param:

```go
//	@Param	payload	body	CreateUserRequest	true	"User payload"
```

Where `CreateUserRequest` is a struct defined in the `handlers` package (or aliased in `swagger_types.go`).

## Tags

Group related endpoints with the same `@Tags` value. Current tags:

| Tag | Endpoints |
|-----|-----------|
| `general` | `GET /` |
| `ops` | `GET /health` |

Add new tags as needed — swag collects them automatically.
