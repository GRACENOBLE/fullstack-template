---
topic: data-fetching
last_verified: 2026-06-25
sources:
  - app/page.tsx
  - app/layout.tsx
  - server/trpc.ts
  - server/routers/_app.ts
  - server/routers/health.ts
  - lib/trpc/client.tsx
  - lib/trpc/server.ts
  - lib/trpc/utils.ts
  - app/providers.tsx
  - lib/user-profile.ts
  - app/(dashboard)/dashboard/page.tsx
  - app/(dashboard)/dashboard/ProfileCard.tsx
---

# Data Fetching

## Core principle
Fetch data in Server Components — not in client components, not via `useEffect`.
Server Components run only on the server, so `fetch` calls go directly to the backend without CORS restrictions or exposed API keys.

## Server Component fetch (preferred pattern)
```tsx
// app/users/page.tsx — Server Component (no "use client")
export default async function UsersPage() {
  const res = await fetch('http://localhost:8080/api/users', {
    // Next.js 16: opt into caching explicitly
    next: { revalidate: 60 }, // revalidate every 60s
    // or: cache: 'no-store'  // always fresh
  });

  if (!res.ok) throw new Error('Failed to fetch users');
  const users = await res.json();

  return <ul>{users.map(u => <li key={u.id}>{u.name}</li>)}</ul>;
}
```

## Server Actions (mutations)
For form submissions and data mutations, use Server Actions — not client-side `fetch`:

```tsx
// app/users/actions.ts
'use server';

export async function createUser(formData: FormData) {
  const name = formData.get('name') as string;
  await fetch('http://localhost:8080/api/users', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name }),
  });
  revalidatePath('/users');
}
```

```tsx
// app/users/page.tsx — use the action
import { createUser } from './actions';
<form action={createUser}>
  <input name="name" />
  <button type="submit">Add</button>
</form>
```

## When client-side fetching is acceptable
Only when data depends on runtime browser state (e.g., user interaction, live updates).
Use tRPC with React Query (via `trpc.<router>.<procedure>.useQuery()`) for client-side data — never bare `useEffect + fetch`.

## Backend URL
Backend runs at `http://localhost:8080` in development.
Store in an env var for production: `process.env.NEXT_PUBLIC_API_URL` (client-accessible) or `process.env.API_URL` (server-only).

## Error handling
- Throw in `async` Server Components — Next.js catches it and renders `error.tsx`.
- Always check `res.ok` before `.json()`.
- Use `notFound()` from `next/navigation` for 404 cases.

## Caching (Next.js 16)
Next.js 16 changes default caching behavior from v14. Check `web/node_modules/next/dist/docs/` for the current defaults before assuming cached or uncached behavior.

## tRPC (preferred for typed API calls)

For calls to the Go backend that benefit from end-to-end type safety, use tRPC instead of bare `fetch`.

**When to use tRPC vs. plain `fetch`:**
- Use tRPC when calling procedures already defined in `server/routers/` — you get compile-time type inference and no manual `res.json()` casting.
- Use plain `fetch` for one-off external APIs, webhooks, or cases where a tRPC router would be disproportionate overhead.

**Server Components** — use `createServerCaller` from `lib/trpc/server.ts`:
```tsx
import { createServerCaller } from '@/lib/trpc/server'

export default async function HealthPage() {
  const caller = await createServerCaller()
  const health = await caller.health.query()
  return <p>Status: {health.status}</p>
}
```

**Client Components** — use `trpc.<router>.<procedure>.useQuery()` or `useMutation()` from `lib/trpc/client.tsx`:
```tsx
'use client'
import { trpc } from '@/lib/trpc/client'

export function HealthStatus() {
  const { data, isLoading, isError } = trpc.health.query.useQuery()
  if (isLoading) return <p>Loading…</p>
  if (isError) return <p>Error</p>
  return <p>Status: {data?.status}</p>
}
```

See [`docs/trpc.md`](trpc.md) for full setup details, context, protected procedures, and how to add new routers.

## Backend fetch utility pattern (`lib/user-profile.ts`)

For fetches to the Go backend that don't warrant a tRPC procedure, place a typed utility function in `lib/`. The `fetchUserProfile` utility in `lib/user-profile.ts` demonstrates the canonical shape:

```ts
// lib/user-profile.ts
export interface UserProfile {
  uid: string
  email: string
  displayName: string
}

export async function fetchUserProfile(userId: string): Promise<UserProfile> {
  const backendUrl = process.env.BACKEND_URL   // server-only env var
  if (!backendUrl) throw new Error('BACKEND_URL environment variable is not set')

  const res = await fetch(`${backendUrl}/api/v1/me`, {
    cache: 'no-store',
    headers: { 'X-User-Id': userId },
  })

  if (!res.ok) {
    throw new Error(`Failed to fetch user profile: ${res.status} ${res.statusText}`)
  }

  const body = (await res.json()) as { data: UserProfile }
  return body.data   // unwrap Go backend { "data": ... } envelope
}
```

Key rules:
- Use `process.env.BACKEND_URL` (no `NEXT_PUBLIC_` prefix) — the utility runs server-side only.
- Always assert `res.ok` before calling `.res.json()`.
- Unwrap the Go backend's `{ "data": ... }` envelope before returning.
- `cache: 'no-store'` — user-specific data must not be cached across requests.

### Server Component with fallback

`app/(dashboard)/dashboard/page.tsx` calls `fetchUserProfile` inside a Server Component and falls back to session data when the backend is unreachable:

```tsx
// app/(dashboard)/dashboard/page.tsx
export default async function DashboardPage() {
  const session = await auth()
  if (!session) redirect('/login')

  const fallbackProfile: UserProfile = {
    uid: session.user?.id ?? '',
    email: session.user?.email ?? '',
    displayName: session.user?.name ?? session.user?.email ?? 'User',
  }

  let profile: UserProfile = fallbackProfile
  try {
    profile = await fetchUserProfile(fallbackProfile.uid)
  } catch {
    profile = fallbackProfile
  }

  return <ProfileCard profile={profile} />
}
```

`ProfileCard` (`app/(dashboard)/dashboard/ProfileCard.tsx`) is a pure Server Component — no `"use client"` directive — that receives a `UserProfile` prop and renders it.

> Auth note: NextAuth is configured with the JWT/Credentials strategy. The original Firebase ID token is not stored in the session — only decoded claims (uid, email, name) are. `fetchUserProfile` therefore sends `X-User-Id` instead of a Firebase Bearer token.
