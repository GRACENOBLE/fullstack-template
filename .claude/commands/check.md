Run a full quality gate across both backend and frontend. Report pass/fail for each step. Stop on the first critical failure.

## Backend
```bash
cd backend && go vet ./...
cd backend && make test
```

## Frontend
```bash
cd frontend && pnpm lint
cd frontend && pnpm build
```

Show all output for failing steps. For each failure, identify the root cause and suggest a fix — but do not apply fixes without confirmation.

End with a summary table:

| Check | Status |
|---|---|
| go vet | ✓ / ✗ |
| backend tests | ✓ / ✗ |
| pnpm lint | ✓ / ✗ |
| pnpm build | ✓ / ✗ |
