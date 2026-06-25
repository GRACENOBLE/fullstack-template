# ADR 0006 — Redis + Asynq for background jobs and event streaming

**Status:** Accepted  
**Date:** 2026-06-25

## Context

The template needs infrastructure for:
- Deferred/background tasks (welcome emails, notification dispatch, data processing)
- Event-driven messaging between services or components
- Rate-limited or scheduled work (cron jobs, retries with backoff)

Candidates evaluated: Asynq + Redis, RabbitMQ, SQS, Temporal.

## Decision

Use **Redis** (`github.com/redis/go-redis/v9`) as the backing store with two layers:

1. **Asynq** (`github.com/hibiken/asynq`) for task queues — priority queues, retries, scheduled tasks, and the Asynqmon web UI for monitoring.
2. **Redis Streams** (built-in Redis) for event sourcing / pub-sub patterns between producers and consumers.

**Both features are opt-in.** Omitting `REDIS_URL` from the environment disables both the cache and all queue/stream features. The Redis Streams consumer is intentionally **not auto-started** — the wiring block in `cmd/api/main.go` is commented out and documented in `backend/docs/streams.md`.

## Consequences

### Positive
- Asynq provides exactly-once delivery guarantees via Redis SET NX, plus automatic retries and dead-letter queues.
- The Asynqmon UI (available at `/admin/queues` in debug mode) gives visibility into queues without a separate tool.
- Redis Streams' consumer group model allows multiple workers to process events in parallel with acknowledgement semantics.
- A single Redis instance serves the cache (`redis.Client`), Asynq broker, and stream producer — reducing infrastructure complexity.

### Negative / trade-offs
- Redis adds an operational dependency. If Redis is unavailable, all queue and stream features fail.
- Asynq does not support message ordering within a single queue; use Redis Streams if strict ordering is required.
- The Redis Streams consumer is opt-in by design — contributors must uncomment and adapt the wiring in `cmd/api/main.go` for their use case.
- Redis Streams have no built-in dead-letter queue; consumer-side retry logic must be implemented manually.
