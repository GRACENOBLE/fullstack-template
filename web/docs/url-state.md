---
topic: url-state
last_verified: 2026-06-27
sources:
  - app/providers.tsx
---

# URL Search-Param State (nuqs)

[nuqs](https://nuqs.47ng.com) synchronises React state with URL search params. It replaces manual `useSearchParams` + `router.push` wiring and works with Next.js App Router out of the box.

## Adapter

`NuqsAdapter` is already mounted as the outermost provider in `app/providers.tsx`. No extra setup is needed.

## Basic usage

```tsx
'use client'
import { useQueryState, parseAsInteger } from 'nuqs'

export function PageSizeSelector() {
  const [pageSize, setPageSize] = useQueryState('pageSize', parseAsInteger.withDefault(20))

  return (
    <select value={pageSize} onChange={e => setPageSize(Number(e.target.value))}>
      <option value={10}>10</option>
      <option value={20}>20</option>
      <option value={50}>50</option>
    </select>
  )
}
```

The URL becomes `?pageSize=50`; removing the param resets to the default.

## Multiple params at once

```tsx
'use client'
import { useQueryStates, parseAsString, parseAsInteger } from 'nuqs'

const [{ q, page }, setSearch] = useQueryStates({
  q: parseAsString.withDefault(''),
  page: parseAsInteger.withDefault(1),
})

// Merge-update (only touches specified keys):
setSearch({ page: 2 })

// Full replace:
setSearch({ q: 'hello', page: 1 })
```

## Built-in parsers

| Parser | URL value | JS value |
|---|---|---|
| `parseAsString` | `?q=hello` | `'hello'` |
| `parseAsInteger` | `?page=2` | `2` |
| `parseAsBoolean` | `?open=true` | `true` |
| `parseAsFloat` | `?price=9.99` | `9.99` |
| `parseAsIsoDateTime` | `?date=2026-06-27T...` | `Date` |
| `parseAsArrayOf(parseAsString)` | `?tags=a&tags=b` | `['a', 'b']` |
| `parseAsJson(zodSchema.parse)` | `?filter=%7B...%7D` | parsed object |

## Server Component access

Server Components receive search params as a prop — no hook needed:

```tsx
// app/search/page.tsx
export default async function SearchPage({
  searchParams,
}: {
  searchParams: Promise<{ q?: string; page?: string }>
}) {
  const { q = '', page = '1' } = await searchParams
  // pass to server action / tRPC caller
}
```

Use `createSearchParamsCache` from nuqs/server to parse and validate in one step:

```tsx
import { createSearchParamsCache, parseAsString, parseAsInteger } from 'nuqs/server'

export const searchParamsCache = createSearchParamsCache({
  q: parseAsString.withDefault(''),
  page: parseAsInteger.withDefault(1),
})

export default async function SearchPage({ searchParams }: { searchParams: Promise<Record<string, string>> }) {
  const { q, page } = searchParamsCache.parse(await searchParams)
  // ...
}
```

## When to use nuqs vs. plain `useSearchParams`

| Scenario | Use |
|---|---|
| Syncing filter/sort/pagination UI with the URL | `useQueryState` / `useQueryStates` |
| Reading params in a Server Component | `searchParams` prop (or `createSearchParamsCache`) |
| One-off read-only access in a Client Component | `useSearchParams()` (Next.js built-in) |
| Programmatic navigation without state binding | `router.push` / `router.replace` |

Prefer nuqs whenever a Client Component both reads and writes a search param — it handles serialisation, defaults, and shallow routing automatically.
