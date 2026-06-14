---
topic: routing
last_verified: 2026-06-14
sources:
  - app/layout.tsx
  - app/page.tsx
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

## Root layout (`app/layout.tsx`)
- Must export `metadata` and a default `RootLayout` component.
- Sets up fonts (Geist Sans + Geist Mono via `next/font/google`), global CSS, `<html>` and `<body>`.
- Font CSS variables: `--font-geist-sans`, `--font-geist-mono` — used in `globals.css` via `@theme inline`.
- `<html>` carries font variable classes; `<body>` has `min-h-full flex flex-col`.

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
Use `(groupName)/` to group routes without affecting the URL:
```
app/
  (marketing)/
    about/page.tsx  → /about
  (app)/
    dashboard/page.tsx → /dashboard
```

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
