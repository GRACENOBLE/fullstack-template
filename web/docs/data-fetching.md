---
topic: data-fetching
last_verified: 2026-06-14
sources:
  - app/page.tsx
  - app/layout.tsx
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
Use `SWR` or `TanStack Query` for client-side data — never bare `useEffect + fetch`.

## Backend URL
Backend runs at `http://localhost:8080` in development.
Store in an env var for production: `process.env.NEXT_PUBLIC_API_URL` (client-accessible) or `process.env.API_URL` (server-only).

## Error handling
- Throw in `async` Server Components — Next.js catches it and renders `error.tsx`.
- Always check `res.ok` before `.json()`.
- Use `notFound()` from `next/navigation` for 404 cases.

## Caching (Next.js 16)
Next.js 16 changes default caching behavior from v14. Check `web/node_modules/next/dist/docs/` for the current defaults before assuming cached or uncached behavior.
