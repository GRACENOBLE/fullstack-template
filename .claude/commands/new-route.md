Scaffold a new API route end-to-end in the Go backend.

You need from the user:
- Route path (e.g. `/api/users`)
- HTTP method (GET / POST / PUT / DELETE / PATCH)
- Brief description of what it does
- Whether it needs database access

## Steps to follow exactly

### 1. Register the route in `backend/internal/server/routes.go`
Add inside `RegisterRoutes()`, grouped with related routes:
```go
r.GET("/api/users", s.listUsersHandler)
```

### 2. Add the handler in `backend/internal/server/routes.go`
Follow this exact signature:
```go
func (s *Server) listUsersHandler(c *gin.Context) {
    // handler body
    c.JSON(http.StatusOK, gin.H{"data": result})
}
```

### 3. If the route needs DB access
In `backend/internal/database/database.go`:
- Add the method signature to the `Service` interface
- Implement it on `*service` using `s.db.QueryContext` / `s.db.ExecContext`
- Use parameterized queries only (never string-concatenated SQL)

### 4. Add an integration test in `backend/internal/database/database_test.go`
Follow the `TestMain` + `mustStartPostgresContainer()` pattern already in that file.
Use table-driven tests when testing multiple cases.

## After scaffolding
Run:
```bash
cd backend && go vet ./...
cd backend && make test
```
Fix any errors before finishing.
