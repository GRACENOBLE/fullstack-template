---
topic: auth
last_verified: 2026-06-24
sources:
  - auth.ts
  - lib/firebase.ts
  - lib/firebase-admin.ts
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

Firebase-first auth with NextAuth v5 (Auth.js) for session management.

**Flow overview:**
1. Client signs in with Firebase Auth (Google popup or email/password).
2. Client retrieves a Firebase ID token (`user.getIdToken()`).
3. Client calls NextAuth's Credentials provider with `{ idToken }`.
4. `auth.ts` verifies the token with Firebase Admin SDK (`verifyIdToken`), extracts claims, and returns a NextAuth user.
5. NextAuth issues a JWT session cookie.

## Packages

```json
"next-auth": "beta",
"firebase": "^12.x",
"firebase-admin": "^14.x",
"react-hook-form": "^7.x",
"@hookform/resolvers": "^5.x"
```

## Config (`auth.ts`)

Lives at the web root. Exports four named values used throughout the app:

```ts
export const { handlers, auth, signIn, signOut } = NextAuth({ ... })
```

- **`handlers`** — `{ GET, POST }` for the route handler.
- **`auth()`** — async function; call in Server Components, Server Actions, and proxy to read the session.
- **`signIn(provider, options)`** — trigger sign-in from server context.
- **`signOut(options)`** — trigger sign-out from server context.

### Provider

Single **Credentials** provider that accepts `{ idToken: string }`.

`authorize()`:
1. Validates `idToken` is a non-empty string with Zod.
2. Calls `verifyFirebaseToken(idToken)` (Firebase Admin SDK) — verifies signature, expiry, audience, and issuer.
3. Returns `{ id, email, name, image }` from the decoded claims, or `null` on any failure.

### Session strategy

```ts
session: { strategy: 'jwt' }
```

### Custom sign-in page

```ts
pages: { signIn: '/login' }
```

## Firebase Admin SDK (`lib/firebase-admin.ts`)

Initialises the Admin app once (singleton) and exports `verifyFirebaseToken`:

```ts
export async function verifyFirebaseToken(idToken: string) {
  return getAuth(getAdminApp()).verifyIdToken(idToken)
}
```

Initialization precedence:
1. `FIREBASE_SERVICE_ACCOUNT_JSON` (full service account JSON string) — works everywhere.
2. `FIREBASE_PROJECT_ID` alone — works on GCP with Application Default Credentials.

## Firebase Client SDK (`lib/firebase.ts`)

Exports `getFirebaseAuth()` for use in client components. Initialises the Firebase app once using `NEXT_PUBLIC_FIREBASE_*` env vars.

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
    loginUrl.searchParams.set("callbackUrl", req.nextUrl.pathname + req.nextUrl.search)
    return NextResponse.redirect(loginUrl)
  }
})

export const config = {
  matcher: ["/dashboard/:path*", "/settings/:path*"],
}
```

Unauthenticated requests to `/dashboard/*` or `/settings/*` are redirected to `/login?callbackUrl=<original-path-and-query>`.

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
      <Toaster richColors position="top-right" />
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

- **`loginSchema`** — `{ email: z.email(), password: string (min 1) }`.
- **`registerSchema`** — `{ name: string (min 2), email: z.email(), password: string (min 8), confirmPassword: string }` with `.refine` for password match.
- Exports inferred types `LoginFormValues` and `RegisterFormValues`.

Used with `standardSchemaResolver` from `@hookform/resolvers/standard-schema` (required for Zod v4).

## Auth components (`features/auth/components/`)

All are `'use client'` components. Errors are surfaced via Sonner toasts, not inline error state.

### `LoginForm`

`react-hook-form` with `standardSchemaResolver(loginSchema)`. On submit:
1. `signInWithEmailAndPassword(getFirebaseAuth(), email, password)` via Firebase.
2. Retrieves ID token from the credential.
3. `signIn('credentials', { idToken, redirect: false })` via NextAuth.
4. Pushes to `/dashboard` on success; shows `toast.error` on failure.

### `RegisterForm`

Same form setup with `registerSchema`. On submit:
1. `createUserWithEmailAndPassword(getFirebaseAuth(), email, password)`.
2. `updateProfile(user, { displayName: name })`.
3. Retrieves ID token.
4. `signIn('credentials', { idToken, redirect: false })`.
5. Pushes to `/dashboard` on success.
Handles `auth/email-already-in-use` with a specific toast message.

### `GoogleSignInButton`

On click:
1. `signInWithPopup(getFirebaseAuth(), new GoogleAuthProvider())`.
2. Retrieves ID token.
3. `signIn('credentials', { idToken, redirect: false })`.
4. Pushes to `/dashboard` on success.
Popup-dismissed errors (`auth/popup-closed-by-user`, `auth/cancelled-popup-request`) are silently ignored. Button is disabled while in-flight.

### `UserMenu`

Reads `{ user, isAuthenticated }` from `useSession()`. Returns `null` when not authenticated. Renders an `Avatar` inside a `DropdownMenu` showing the user's name, email, and a "Sign out" item that calls `useSignOut().signOut`.

## Auth hooks (`features/auth/hooks/`)

### `useSession`

Wraps `useSession` from `next-auth/react`. Returns `UseSessionReturn`:

```ts
{ user: AuthUser | null, isAuthenticated: boolean, isLoading: boolean }
```

### `useSignOut`

Signs out of both NextAuth and Firebase in parallel, then navigates to `/login`:

```ts
await Promise.all([
  nextAuthSignOut({ redirect: false }),
  firebaseSignOut(getFirebaseAuth()),
])
router.push('/login')
```

## Pages

| Route | File | Type |
|---|---|---|
| `/login` | `app/(auth)/login/page.tsx` | Server Component |
| `/register` | `app/(auth)/register/page.tsx` | Server Component |
| `/dashboard` | `app/(dashboard)/dashboard/page.tsx` | Server Component |
| `/settings` | `app/(dashboard)/settings/page.tsx` | Server Component |

Login and Register pages are layout-less Server Components that render a centered `Card` containing `GoogleSignInButton`, a divider, the form, and a link to the other page.

## Environment variables

| Variable | Description |
|---|---|
| `AUTH_SECRET` | Secret used to sign JWTs. Generate with `openssl rand -base64 32`. |
| `FIREBASE_PROJECT_ID` | Firebase project ID (server-side). Used when `FIREBASE_SERVICE_ACCOUNT_JSON` is absent and ADC is available (GCP). |
| `FIREBASE_SERVICE_ACCOUNT_JSON` | Full service account JSON as a single-line string. Required outside GCP for `verifyIdToken`. |
| `NEXT_PUBLIC_FIREBASE_API_KEY` | Firebase web API key. |
| `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN` | Firebase auth domain. |
| `NEXT_PUBLIC_FIREBASE_PROJECT_ID` | Firebase project ID (client-side). |
| `NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET` | Firebase storage bucket. |
| `NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID` | Firebase messaging sender ID. |
| `NEXT_PUBLIC_FIREBASE_APP_ID` | Firebase app ID. |

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
