---
topic: DataTable component (TanStack Table v8)
last_verified: 2026-06-25
sources:
  - components/data-table/DataTable.tsx
  - components/data-table/columns.ts
  - components/data-table/index.ts
  - app/demo/page.tsx
  - app/demo/DemoTable.tsx
  - components/data-table/__tests__/DataTable.test.tsx
---

# DataTable component

## When to use

Use `DataTable` for any list that needs client-side sorting, column filtering, and pagination. It is a Client Component (`"use client"`) backed by TanStack Table v8 (`@tanstack/react-table`).

## Props interface

```tsx
interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
  filterColumn?: string
  filterPlaceholder?: string
}
```

| Prop | Type | Default | Purpose |
|---|---|---|---|
| `columns` | `ColumnDef<TData, TValue>[]` | — | Column definitions (see below) |
| `data` | `TData[]` | — | Row data array |
| `filterColumn` | `string \| undefined` | `undefined` | When set, renders a text input that filters rows by this column's accessor key |
| `filterPlaceholder` | `string` | `"Filter…"` | Placeholder and `aria-label` for the filter input |

Pagination defaults to **10 rows per page** (`initialState: { pagination: { pageSize: 10 } }`). The page size is not currently configurable via props.

## Defining columns

Columns are `ColumnDef<TData, TValue>[]` objects from `@tanstack/react-table`. Place column arrays in `components/data-table/columns.ts` (or a feature-specific file).

```ts
import { type ColumnDef } from "@tanstack/react-table"

export interface User {
  id: string
  name: string
  email: string
  createdAt: string
}

export const userColumns: ColumnDef<User, string>[] = [
  {
    accessorKey: "id",
    header: "ID",
    enableSorting: false,
  },
  {
    accessorKey: "name",
    header: "Name",
    enableSorting: true,
  },
  {
    accessorKey: "createdAt",
    header: "Created At",
    enableSorting: false,
    cell: ({ getValue }) => {
      const raw = getValue()
      const date = new Date(raw)
      if (isNaN(date.getTime())) return raw
      const yyyy = date.getUTCFullYear()
      const mm = String(date.getUTCMonth() + 1).padStart(2, "0")
      const dd = String(date.getUTCDate()).padStart(2, "0")
      return `${yyyy}-${mm}-${dd}`
    },
  },
]
```

- `accessorKey` — maps to a field on `TData`.
- `enableSorting: true` — renders a sort toggle (↕ / ↑ / ↓) in the header; `false` disables it.
- `cell` — optional custom renderer; receives `{ getValue }` and must return a `React.ReactNode` (or a primitive).

## Using the component

```tsx
import { DataTable } from "@/components/data-table"
import { userColumns } from "@/components/data-table"

export function DemoTable() {
  return (
    <DataTable
      columns={userColumns}
      data={users}
      filterColumn="name"
      filterPlaceholder="Search by name…"
    />
  )
}
```

`DemoTable` must be a Client Component (`"use client"`) because it passes state-affecting columns. The Server Component page (`app/demo/page.tsx`) imports `DemoTable` and renders it inside a static `<main>` layout.

## Demo page

A working demo is available at `/demo` (`app/demo/page.tsx` + `app/demo/DemoTable.tsx`). It renders 25 hardcoded `User` rows with `filterColumn="name"` and shows 10 rows per page.

## Exports

`components/data-table/index.ts` re-exports:

```ts
export { DataTable } from "./DataTable"
export { userColumns, type User } from "./columns"
```

Import from `@/components/data-table` to pick up both.

## Internals

`DataTable` wires four TanStack Table row models:

| Row model | Hook |
|---|---|
| Core | `getCoreRowModel()` |
| Sorting | `getSortedRowModel()` |
| Column filtering | `getFilteredRowModel()` |
| Pagination | `getPaginationRowModel()` |

Sortable headers receive `onClick` → `getToggleSortingHandler()` and `aria-sort` attributes (`"ascending"` / `"descending"` / `"none"`). The filter input calls `column.setFilterValue()` on the `onChange` event. Pagination is controlled by `table.previousPage()` / `table.nextPage()` buttons with `aria-label="Previous page"` / `aria-label="Next page"`.

## Testing

Tests live in `components/data-table/__tests__/DataTable.test.tsx` and run with `pnpm test` (Vitest + `@testing-library/react`, jsdom environment).

Coverage:

| Test | What is asserted |
|---|---|
| Initial render | 10 rows visible for a 25-item dataset |
| Filtering | Typing in the filter input narrows rows to matching entries |
| Pagination — page 1 | Previous button is disabled |
| Pagination — advance | Next click shows page 2; Previous becomes enabled |
| Pagination — last page | Correct row count on the final page |
| Sorting ascending | Clicking a sortable header once sorts that column A→Z |
| Sorting descending | Clicking the same header a second time reverses the sort |

Pattern for selecting rows (avoids picking up the header row):

```tsx
const tbody = screen.getByRole("table").querySelector("tbody")!
const rows = within(tbody).getAllByRole("row")
```
