# Template Readiness Status

Tracks all identified gaps from the June 2026 template analysis. Issues live in the [fullstack-template GitHub project](https://github.com/users/GRACENOBLE/projects/8).

---

## Priority: Critical — Day 1 blockers

| # | Issue | Status |
|---|-------|--------|
| [#44](https://github.com/GRACENOBLE/fullstack-template/issues/44) | `backend` CORS origin hardcoded to localhost:3000 → make config-driven | ✅ Done (merged) |
| [#45](https://github.com/GRACENOBLE/fullstack-template/issues/45) | `web` `.env.example` missing BACKEND_URL, SENTRY_ORG, SENTRY_PROJECT | ✅ Done (merged) |
| [#46](https://github.com/GRACENOBLE/fullstack-template/issues/46) | `backend` Standardized API response envelope (`JSON`, `JSONError` helpers) | ✅ Done (merged) |
| [#47](https://github.com/GRACENOBLE/fullstack-template/issues/47) | `backend` Request ID middleware (`X-Request-ID`, propagated to logs) | ✅ Done (merged) |
| [#48](https://github.com/GRACENOBLE/fullstack-template/issues/48) | `web` Error pages: `not-found.tsx`, `error.tsx`, `global-error.tsx` | ✅ Done (merged) |
| [#49](https://github.com/GRACENOBLE/fullstack-template/issues/49) | `mobile` HTTP client for backend API calls (OkHttp + Firebase token interceptor) | ✅ Done (merged) |

---

## Priority: High — First-week friction

| # | Issue | Status |
|---|-------|--------|
| [#50](https://github.com/GRACENOBLE/fullstack-template/issues/50) | `infra` Root dev script to start all services with one command | ✅ Done (merged) |
| [#51](https://github.com/GRACENOBLE/fullstack-template/issues/51) | `backend` Add `.golangci.yml` linter config | ✅ Done (merged) |
| [#52](https://github.com/GRACENOBLE/fullstack-template/issues/52) | `mobile` Add ktlint and integrate into CI | ✅ Done (merged) |
| [#53](https://github.com/GRACENOBLE/fullstack-template/issues/53) | `infra` Renovate / Dependabot for automated dependency updates | ✅ Done (merged) |
| [#54](https://github.com/GRACENOBLE/fullstack-template/issues/54) | `infra` First-run setup script for new contributors | ✅ Done (merged) |
| [#55](https://github.com/GRACENOBLE/fullstack-template/issues/55) | `backend` Swagger generation check in CI (fail if stale) | ✅ Done (merged) |
| [#56](https://github.com/GRACENOBLE/fullstack-template/issues/56) | `mobile` Loading state and skeleton screen pattern | ✅ Done (merged) |
| [#57](https://github.com/GRACENOBLE/fullstack-template/issues/57) | `mobile` Error state and retry UI pattern (`UiState<T>` sealed class) | ✅ Done (merged) |
| [#58](https://github.com/GRACENOBLE/fullstack-template/issues/58) | `web` Data table with sorting, filtering, and pagination (TanStack Table) | ✅ Done (merged) |

---

## Priority: Medium — Polish

| # | Issue | Status |
|---|-------|--------|
| [#59](https://github.com/GRACENOBLE/fullstack-template/issues/59) | `mobile` Settings screen: show user profile + sign-out button | ✅ Done (merged) |
| [#60](https://github.com/GRACENOBLE/fullstack-template/issues/60) | `web` Dashboard page: fetch and display `/api/v1/me` | ✅ Done (merged) |
| [#61](https://github.com/GRACENOBLE/fullstack-template/issues/61) | `backend` Redis stream consumers: wire with feature flag or document as opt-in | ✅ Done (merged) |
| [#62](https://github.com/GRACENOBLE/fullstack-template/issues/62) | `backend` pprof endpoints for runtime profiling (gated to internal network) | ✅ Done (merged) |

---

## Priority: Low — Nice-to-have

| # | Issue | Status |
|---|-------|--------|
| [#63](https://github.com/GRACENOBLE/fullstack-template/issues/63) | `infra` Architecture Decision Records (ADRs) for key technology choices | 🔁 In review (PR #68) |
| [#64](https://github.com/GRACENOBLE/fullstack-template/issues/64) | `infra` Deployment runbook for staging and production | 🔁 In review (PR #68) |

---

## Legend

| Symbol | Meaning |
|--------|---------|
| ✅ Done | Merged to main |
| 🔁 In review | PR open, pending merge |
| ⬜ Open | Not started |
| 🚧 In progress | Branch exists, work ongoing |

---

_Last updated: 2026-06-25 — PR #67 merged (medium priority #59–#62); PR #68 open (nice-to-haves #63–#64)._
