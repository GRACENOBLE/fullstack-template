# ADR 0003 — Next.js 16 App Router + React 19

**Status:** Accepted  
**Date:** 2026-06-25

## Context

The template needs a web layer that:
- Supports server-side rendering and static generation without a separate server
- Enables Server Components to reduce client JavaScript bundle size
- Has a mature deployment target (Vercel, Railway, Docker)
- Works with TypeScript 5 and Tailwind CSS v4

Candidates evaluated: Next.js App Router, Remix, SvelteKit.

## Decision

Use **Next.js 16** with the **App Router** (not the legacy Pages Router), **React 19**, and **TypeScript 5**.

Key conventions enforced:
- Default to **Server Components** — add `"use client"` only when browser APIs or React hooks are needed, pushed to the smallest possible component.
- Route groups `(auth)/` and `(dashboard)/` organise unauthenticated vs. protected pages without affecting URLs.
- **Tailwind CSS v4** for all styling — no CSS modules, no inline styles.
- **tRPC + TanStack Query** for type-safe client→server RPC calls.
- Error boundaries: `app/error.tsx`, `app/not-found.tsx`, `app/global-error.tsx`.

## Consequences

### Positive
- Server Components eliminate round-trips for data that doesn't need to be interactive.
- App Router colocates layouts, loading states, and error boundaries with the routes they belong to.
- TypeScript end-to-end (tRPC infers types from router definitions, no manual API typing).
- Vercel deploys with zero config; Railway and Docker are also viable.

### Negative / trade-offs
- The `"use client"` boundary requires discipline — violations silently increase the bundle.
- App Router's caching model (fetch cache, `no-store`, `revalidate`) is nuanced; misuse causes stale data or redundant fetches.
- Server Component testing requires Playwright E2E; unit testing is limited to extracted logic and Client Components.
- Next.js 16 is a major version with breaking changes from Next.js 13/14; contributors should read `node_modules/next/dist/docs/` before writing code.
