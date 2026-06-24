---
topic: auth
last_verified: 2026-06-24
sources:
  - auth.ts
  - app/api/auth/[...nextauth]/route.ts
  - proxy.ts
  - server/trpc.ts
  - server/routers/auth.ts
  - app/providers.tsx
  - features/auth/types.ts
  - features/auth/validation.ts
  - features/auth/components/LoginForm.tsx
  - features/auth/components/RegisterForm.tsx
  - features/auth/components/GoogleSignInButton.tsx
  - features/auth/components/UserMenu.tsx
  - features/auth/hooks/useSession.ts
  - features/auth/hooks/useSignOut.ts
  - app/(auth)/login/page.tsx
  - app/(auth)/register/page.tsx
  - app/(dashboard)/dashboard/page.tsx
---

# Authentication

NextAuth v5 (`next-auth@beta`) with Google OAuth and Credentials providers.

## Packages

```json
"next-auth": "beta",
"react-hook-form": "^7.x",
"@hookform/resolvers": "^3.x"
```

## Config (`auth.ts`)

Lives at the web root. Exports four named values used throughout the app:

```ts
export const { handlers, auth, signIn, signOut } = NextAuth({ ... })
```

- **`handlers`** — `{ GET, POST }` for the route handler.
- **`auth()`** — async function; call in Server Components, Server Actions, and middleware to read the session.
- **`signIn(provider, options)`** — trigger sign-in from server context.
- **`signOut(options)`** — trigger sign-out from server context.

### Providers

**Google** — reads `AUTH_GOOGLE_ID` / `AUTH_GOOGLE_SECRET` from env.

**Credentials** — validates email + password with Zod, then proxies a `POST` to `${NEXT_PUBLIC_API_URL}/auth/login`. Returns `{ id, email, name }` on success or `null` on failure.

### Session strategy

```ts
session: { strategy: 'jwt' }
```

### Custom sign-in page

```ts
pages: { signIn: '/login' }
```

## Route handler (`app/api/auth/[...nextauth]/route.ts`)

```ts
import { handlers } from "@/auth"
export const { GET, POST } = handlers
```

All NextAuth HTTP endpoints (`/api/auth/callback/*`, `/api/auth/session`, etc.) are handled here.

## Route protection

### Proxy (`proxy.ts`)

```ts
import { auth } from "@/auth"

export default auth((req: NextAuthRequest) => {
  if (!req.auth) {
    const loginUrl = new URL("/login", req.url)
    loginUrl.searchParams.set("callbackUrl", req.nextUrl.pathname)
    return NextResponse.redirect(loginUrl)
  }
})

export const config = {
  matcher: ["/dashboard/:path*"],
}
```

Unauthenticated requests to `/dashboard/*` are redirected to `/login?callbackUrl=<original-path>`.

### Additional guard in Server Components

Protected pages call `auth()` directly and redirect if no session is returned:

```ts
// app/(dashboard)/dashboard/page.tsx
const session = await auth()
if (!session) redirect('/login')
```

## Reading the session

### In a Server Component

```ts
import { auth } from '@/auth'

const session = await auth()
// session is Session | null
```

### In a Client Component

Use the `useSession` hook from `features/auth/hooks/useSession`:

```ts
import { useSession } from '@/features/auth/hooks/useSession'

const { user, isAuthenticated, isLoading } = useSession()
```

Returns `{ user: AuthUser | null, isAuthenticated: boolean, isLoading: boolean }`.
Internally wraps `useSession` from `next-auth/react` and casts to the local `AuthSession` type.

## SessionProvider (`app/providers.tsx`)

`<SessionProvider>` from `next-auth/react` wraps `<TRPCProvider>`:

```tsx
export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <SessionProvider>
      <TRPCProvider>{children}</TRPCProvider>
    </SessionProvider>
  )
}
```

Required for `useSession()` to work in Client Components.

## Types (`features/auth/types.ts`)

```ts
export interface AuthUser {
  id?: string | null
  name?: string | null
  email?: string | null
  image?: string | null
}

export interface AuthSession {
  user: AuthUser
  expires: string
}
```

## Validation (`features/auth/validation.ts`)

Zod v4 schemas:

- **`loginSchema`** — `{ email: string, password: string }` (password min 1).
- **`registerSchema`** — `{ name: string (min 2), email: string, password: string (min 8), confirmPassword: string }` with `.refine` for password match.
- Exports inferred types `LoginFormValues` and `RegisterFormValues`.

## Auth components (`features/auth/components/`)

All are `'use client'` components.

### `LoginForm`

`react-hook-form` with `standardSchemaResolver(loginSchema)` from `@hookform/resolvers/standard-schema`. Calls `signIn('credentials', { redirect: false })` from `next-auth/react` on submit, then pushes to `/dashboard` on success. Displays an inline error string on failure.

### `RegisterForm`

Same form setup with `registerSchema`. On submit: POSTs to `${NEXT_PUBLIC_API_URL}/auth/register`, then calls `signIn('credentials', { redirect: false })` to sign in immediately after registration. Pushes to `/dashboard` on success.

### `GoogleSignInButton`

Calls `signIn('google', { callbackUrl: '/dashboard' })` on click. Renders a `Button` with the Google SVG logo inline.

### `UserMenu`

Reads `{ user, isAuthenticated }` from `useSession()`. Returns `null` when not authenticated. Renders an `Avatar` inside a `DropdownMenu` showing the user's name, email, and a "Sign out" item that calls `useSignOut().signOut`.

## Auth hooks (`features/auth/hooks/`)

### `useSession`

Wraps `useSession` from `next-auth/react`. Returns `UseSessionReturn`:

```ts
{ user: AuthUser | null, isAuthenticated: boolean, isLoading: boolean }
```

### `useSignOut`

Wraps `signOut` from `next-auth/react` with `{ redirect: false }`, then calls `router.push('/login')`.

```ts
const { signOut } = useSignOut()
```

## Pages

| Route | File | Type |
|---|---|---|
| `/login` | `app/(auth)/login/page.tsx` | Server Component |
| `/register` | `app/(auth)/register/page.tsx` | Server Component |
| `/dashboard` | `app/(dashboard)/dashboard/page.tsx` | Server Component |

Login and Register pages are layout-less Server Components that render a centered `Card` containing `GoogleSignInButton`, a divider, the form, and a link to the other page.

## Environment variables

| Variable | Description |
|---|---|
| `AUTH_SECRET` | Secret used to sign JWTs. Generate with `openssl rand -base64 32`. |
| `AUTH_GOOGLE_ID` | Google OAuth client ID. |
| `AUTH_GOOGLE_SECRET` | Google OAuth client secret. |
| `NEXT_PUBLIC_API_URL` | Go backend base URL (e.g. `http://localhost:8080`). Used by Credentials provider and RegisterForm. |

## Testing

`createTRPCContext` calls `auth()` internally. Any test that instantiates a context must mock `@/auth`:

```ts
vi.mock('@/auth', () => ({
  auth: vi.fn(),
}))

import { auth } from '@/auth'
const mockAuth = vi.mocked(auth)

// provide a session:
mockAuth.mockResolvedValue(validSession)
// or simulate unauthenticated:
mockAuth.mockResolvedValue(null)
```

See `server/routers/__tests__/auth.test.ts` for the full test pattern.
