---
topic: trpc
last_verified: 2026-06-24
sources:
  - server/trpc.ts
  - server/routers/_app.ts
  - server/routers/health.ts
  - server/routers/auth.ts
  - server/routers/notifications.ts
  - app/api/trpc/[trpc]/route.ts
  - lib/trpc/client.tsx
  - lib/trpc/server.ts
  - app/providers.tsx
  - auth.ts
  - app/api/auth/[...nextauth]/route.ts
---

# tRPC

tRPC v11 with React Query v5, wired into Next.js App Router.

## Package setup

```json
"@trpc/server": "^11.18.0",
"@trpc/client": "^11.18.0",
"@trpc/react-query": "^11.18.0",
"@tanstack/react-query": "^5.101.1",
"zod": "^4.4.3"
```

## Context and initialization (`server/trpc.ts`)

```ts
export interface TRPCContext {
  req: NextRequest
  session: Session | null
}

export async function createTRPCContext({ req }: { req: NextRequest }): Promise<TRPCContext>
```

`createTRPCContext` calls `auth()` from NextAuth (`@/auth`) to resolve the current session, then returns `{ req, session }`.

### Procedure types

**`publicProcedure`** — alias for `t.procedure`. No auth check.

**`protectedProcedure`** — middleware runs before the handler:
- Checks `ctx.session?.user`.
- Throws `TRPCError({ code: 'UNAUTHORIZED' })` if `session` is `null` or `session.user` is absent.

**`createCallerFactory`** — exported from `t.createCallerFactory`; used by `lib/trpc/server.ts` to build server-side callers.

## Routers (`server/routers/`)

### Root router (`server/routers/_app.ts`)

```ts
export const appRouter = router({
  health: healthRouter,
  auth: authRouter,
  notifications: notificationsRouter,
})

export type AppRouter = typeof appRouter
```

`AppRouter` is the single type exported to the client.

### `health` router (`server/routers/health.ts`)

Uses `publicProcedure`. Fetches `GET ${BACKEND_URL}/health` (defaults to `http://localhost:8080`) and returns `{ status: string; database: string }`. Throws a plain `Error` if the backend responds with a non-2xx status.

### `auth` router (`server/routers/auth.ts`)

Uses `protectedProcedure`.
- `session` — query, returns `{ authenticated: true, user: ctx.session?.user ?? null }`.
- `signOut` — mutation, returns `{ success: true }`.

### `notifications` router (`server/routers/notifications.ts`)

Uses `protectedProcedure` with Zod input validation. Current stubs:
- `registerFcmToken` — mutation, input `{ token: z.string().min(1) }`, returns `{ registered: true, token }`.
- `list` — query, returns `[]`.

## HTTP handler (`app/api/trpc/[trpc]/route.ts`)

```ts
const handler = (req: NextRequest) =>
  fetchRequestHandler({
    endpoint: '/api/trpc',
    req,
    router: appRouter,
    createContext: () => createTRPCContext({ req }),
  })

export { handler as GET, handler as POST }
```

All tRPC requests (batch GET and mutation POST) hit `/api/trpc/[trpc]`.

## Client setup (`lib/trpc/client.tsx`)

```ts
export const trpc = createTRPCReact<AppRouter>()
```

`trpc` is the typed client used in Client Components.

### `getBaseUrl()`

- In the browser: returns `''` (relative URL, same origin).
- On Vercel: returns `https://${process.env.VERCEL_URL}`.
- Elsewhere (local server-side): returns `'http://localhost:3000'`.

### `TRPCProvider`

```tsx
export function TRPCProvider({ children }: { children: React.ReactNode })
```

Creates a `QueryClient` and a `trpc` HTTP batch link client, both memoized in `useState`. Wraps children with `trpc.Provider` and `QueryClientProvider`. This is a `'use client'` component.

## Server-side caller (`lib/trpc/server.ts`)

```ts
export const createServerCaller = cache(async () => {
  const headerList = await headers()
  const req = new Request('http://internal', { headers: headerList }) as NextRequest
  const ctx = await createTRPCContext({ req })
  return createCaller(ctx)
})
```

`createServerCaller` is wrapped in React's `cache()` so it is deduplicated per request. It forwards the incoming request headers (including `Cookie` and `Authorization`) to the context, which means `protectedProcedure` checks work for server-rendered pages.

**Usage in a Server Component:**

```tsx
import { createServerCaller } from '@/lib/trpc/server'

export default async function HealthPage() {
  const caller = await createServerCaller()
  const health = await caller.health.query()
  return <p>Backend status: {health.status}</p>
}
```

## Provider wiring (`app/providers.tsx` + `app/layout.tsx`)

`app/providers.tsx` is a thin `'use client'` wrapper. `SessionProvider` (from `next-auth/react`) wraps `TRPCProvider` so both session and React Query contexts are available to all Client Components:

```tsx
'use client'
import { TRPCProvider } from '@/lib/trpc/client'
import { SessionProvider } from 'next-auth/react'

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <SessionProvider>
      <TRPCProvider>{children}</TRPCProvider>
    </SessionProvider>
  )
}
```

`app/layout.tsx` wraps `{children}` with `<Providers>` inside `<body>`:

```tsx
<body className="min-h-full flex flex-col">
  <Providers>{children}</Providers>
</body>
```

## Usage patterns

### Server Component

```tsx
// No 'use client' — runs on the server
import { createServerCaller } from '@/lib/trpc/server'

export default async function StatusPage() {
  const caller = await createServerCaller()
  const health = await caller.health.query()
  return <p>{health.status}</p>
}
```

### Client Component

```tsx
'use client'
import { trpc } from '@/lib/trpc/client'

export function HealthWidget() {
  const { data, isLoading, isError } = trpc.health.query.useQuery()

  if (isLoading) return <p>Loading…</p>
  if (isError) return <p>Unavailable</p>
  return <p>Status: {data.status} / DB: {data.database}</p>
}
```

### Mutation in a Client Component

```tsx
'use client'
import { trpc } from '@/lib/trpc/client'

export function SignOutButton() {
  const signOut = trpc.auth.signOut.useMutation()
  return (
    <button onClick={() => signOut.mutate()} disabled={signOut.isPending}>
      Sign out
    </button>
  )
}
```

## Adding a new router

1. Create `server/routers/<feature>.ts` and export a `router({})` built from `publicProcedure` or `protectedProcedure`.
2. Import it in `server/routers/_app.ts` and add it to `appRouter`.
3. The new procedures are immediately available to `trpc.<feature>.*` in Client Components and `caller.<feature>.*` in Server Components.

```ts
// server/routers/posts.ts
import { publicProcedure, router } from '../trpc'

export const postsRouter = router({
  list: publicProcedure.query(async () => {
    return []
  }),
})
```

```ts
// server/routers/_app.ts
import { postsRouter } from './posts'

export const appRouter = router({
  health: healthRouter,
  auth: authRouter,
  notifications: notificationsRouter,
  posts: postsRouter,   // add here
})
```

## Testing

Procedures are tested directly via `appRouter.createCaller(ctx)` — no HTTP server needed.

`createTRPCContext` calls `auth()` from NextAuth, so **every test file that calls `createTRPCContext` must mock `@/auth`**:

```ts
vi.mock('@/auth', () => ({
  auth: vi.fn(),
}))

import { auth } from '@/auth'
const mockAuth = vi.mocked(auth)
```

Pattern from `server/routers/__tests__/auth.test.ts`:

```ts
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { appRouter } from '../_app'
import { createTRPCContext } from '../../trpc'
import type { Session } from 'next-auth'

vi.mock('@/auth', () => ({ auth: vi.fn() }))

import { auth } from '@/auth'
const mockAuth = vi.mocked(auth)

const validSession: Session = {
  user: { id: '123', email: 'test@example.com', name: 'Test User' },
  expires: '2099-01-01T00:00:00.000Z',
}

function makeContext(session: Session | null = null) {
  const req = new Request('http://localhost/api/trpc') as NextRequest
  mockAuth.mockResolvedValue(session)
  return createTRPCContext({ req })
}

it('throws UNAUTHORIZED when no session present', async () => {
  const ctx = await makeContext(null)
  const caller = appRouter.createCaller(ctx)
  await expect(caller.auth.session()).rejects.toMatchObject({ code: 'UNAUTHORIZED' })
})

it('allows access with valid session', async () => {
  const ctx = await makeContext(validSession)
  const caller = appRouter.createCaller(ctx)
  const result = await caller.auth.session()
  expect(result).toMatchObject({ authenticated: true, user: { email: 'test@example.com' } })
})
```

For procedures that call `fetch`, stub it with `vi.stubGlobal('fetch', mockFetch)` before the test suite.
