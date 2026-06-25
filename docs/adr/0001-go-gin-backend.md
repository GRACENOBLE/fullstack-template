# ADR 0001 — Go 1.25 + Gin for the backend API

**Status:** Accepted  
**Date:** 2026-06-25

## Context

The template needs a backend API that is:
- Statically typed with compile-time safety
- Fast enough to serve real production workloads without tuning
- Simple enough that a new contributor can navigate the codebase on day one
- Compatible with a Clean Architecture layering (domain → usecase → infrastructure → transport)

Candidates evaluated: Go + Gin, Node.js + Fastify, Rust + Axum.

## Decision

Use **Go 1.25** with **Gin v1.12.0** as the HTTP framework.

The Clean Architecture is enforced via Go package boundaries:

```
domain/          ← entities, no external deps
usecase/         ← application logic, depends on domain
infrastructure/  ← DB, cache, third-party integrations
transport/       ← HTTP handlers and middleware
server/          ← wires all layers
cmd/             ← entry point
```

## Consequences

### Positive
- Go's goroutine model handles high concurrency with a tiny memory footprint compared to thread-per-request models.
- Compile-time type checking catches integration errors before they reach production.
- Gin's middleware chain (request ID, rate limiting, CORS, auth) composes cleanly with `gin.HandlerFunc`.
- Single binary deployment — no runtime dependency management.
- `go vet` and golangci-lint enforce style uniformly in CI.

### Negative / trade-offs
- Go's verbosity (explicit error returns, no generics-based magic) means more boilerplate than TypeScript or Python.
- Gin does not support HTTP/2 push or WebSocket natively; Gorilla WebSocket is added for WS support.
- Contributors unfamiliar with Go's `interface`-based dependency injection need a learning curve before they can add new layers cleanly.
