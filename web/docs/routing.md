---
topic: routing
last_verified: 2026-06-25
sources:
  - app/layout.tsx
  - app/page.tsx
  - app/not-found.tsx
  - app/error.tsx
  - app/global-error.tsx
  - app/(auth)/login/page.tsx
  - app/(auth)/register/page.tsx
  - app/(dashboard)/layout.tsx
  - app/(dashboard)/dashboard/page.tsx
  - app/(dashboard)/dashboard/ProfileCard.tsx
  - app/(dashboard)/settings/page.tsx
  - lib/user-profile.ts
  - next.config.ts
---

# Routing

## Router
App Router only. No Pages Router. Never create files in a `pages/` directory.

## File conventions
| File | Purpose |
|---|---|
| `app/layout.tsx` | Shared layout wrapping children. Root layout is required. |
| `app/page.tsx` | Route UI — the public face of a URL segment |
| `app/loading.tsx` | Suspense boundary shown while page data loads |
| `app/error.tsx` | Error boundary for a segment (must be `"use client"`) |
| `app/not-found.tsx` | 404 UI for the segment |
| `app/global-error.tsx` | Top-level error boundary that catches errors in the root layout (must be `"use client"`; must include its own `<html>` and `<body>` tags) |

## Root layout (`app/layout.tsx`)
- Must export `metadata` and a default `RootLayout` component.
- Sets up fonts (Manrope + Geist Mono via `next/font/google`), global CSS, `<html>` and `<body>`.
- Font CSS variables: `--font-sans` (Manrope), `--font-geist-mono` (Geist Mono) — used in `globals.css` via `@theme inline`.
- `<html>` carries font variable classes and `h-full antialiased`; `<body>` has `min-h-full flex flex-col`.
- Children are wrapped in a `<Providers>` component.

## Nested routes
Add route segments as directories under `app/`:
```
app/
  dashboard/
    page.tsx        → /dashboard
    settings/
      page.tsx      → /dashboard/settings
```

## Route groups
Use `(groupName)/` to group routes without affecting the URL. The project currently has two route groups:

```
app/
  (auth)/
    login/page.tsx      → /login    (Server Component; LoginForm + GoogleSignInButton)
    register/page.tsx   → /register (Server Component; RegisterForm + GoogleSignInButton)
  (dashboard)/
    layout.tsx          → shared layout: SidebarProvider + AppSidebar + SidebarInset;
                          reads sidebar_state cookie to set default open state
    dashboard/page.tsx  → /dashboard (auth-guarded; redirects to /login if no session)
    settings/page.tsx   → /settings  (auth-guarded; redirects to /login if no session)
```

Both `(dashboard)` pages call `auth()` from `@/auth` and redirect to `/login` on a missing session.

## Dynamic segments
```
app/users/[id]/page.tsx   → /users/123
```
Access via props: `{ params }: { params: Promise<{ id: string }> }` in Next.js 16.

## Default component type
All route files (`page.tsx`, `layout.tsx`) are **Server Components** by default.
Do not add `"use client"` to layout or page files unless you have a concrete reason.

## Navigation
Use `<Link href="/path">` from `next/link` — never `<a href>` for internal navigation.
Programmatic navigation: `import { useRouter } from 'next/navigation'` (client components only).

## Error pages

Three special files handle runtime errors and missing routes at the root segment level.

| File | Trigger | Component type | Notes |
|---|---|---|---|
| `app/not-found.tsx` | `notFound()` call or unmatched URL | Server Component | Renders a 404 page with a link back to `/`. No `"use client"` directive. |
| `app/error.tsx` | Uncaught error thrown inside a route segment | Client Component | Receives `error: Error & { digest?: string }` and `reset: () => void` props. The "Try again" button calls `reset()` to re-render the segment. |
| `app/global-error.tsx` | Uncaught error in the root layout itself | Client Component | Same props as `error.tsx`. Must render its own `<html>` and `<body>` tags because the root layout is unavailable when this boundary fires. |

`error.tsx` and `global-error.tsx` must have `"use client"` at the top — Next.js requires error boundaries to be Client Components. `not-found.tsx` has no such requirement and is a Server Component.

`error.digest` is an opaque server-generated hash surfaced in both the UI and server logs; display it as a reference string when present so users can report it.
