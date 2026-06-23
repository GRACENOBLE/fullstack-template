<!-- BEGIN:nextjs-agent-rules -->
# This is NOT the Next.js you know

This version has breaking changes — APIs, conventions, and file structure may all differ from your training data. Read the relevant guide in `node_modules/next/dist/docs/` before writing any code. Heed deprecation notices.
<!-- END:nextjs-agent-rules -->

---

## Commands

```bash
pnpm install      # install dependencies
pnpm dev          # dev server → http://localhost:3000
pnpm build        # production build (also runs TypeScript check)
pnpm start        # serve production build
pnpm lint         # ESLint
pnpm test         # Vitest — unit + component tests
pnpm test:watch   # Vitest watch mode (use during TDD)
pnpm test:ui      # Vitest browser UI
```

---

## Project structure

```
app/                        # Next.js App Router — routes only (no business logic)
  layout.tsx                # root layout: Geist fonts, global CSS, <html>/<body>
  page.tsx                  # home page (Server Component)
  globals.css               # Tailwind v4 import + CSS variable definitions
components/
  layout/                   # AppShell, Sidebar, Navbar, Footer
  common/                   # reusable cross-feature UI (Avatar, Badge, Spinner…)
  home/                     # landing / home page components
  ui/                       # shadcn primitives — populated by CLI, do not hand-create
features/                   # one subfolder per domain feature (auth/, notifications/…)
  <feature>/
    components/             # feature-scoped UI
    hooks/                  # feature-scoped React hooks
    types.ts                # feature-scoped TypeScript types
lib/                        # pure utilities and non-React helpers
public/                     # static assets
server/                     # tRPC routers and server-side logic
types/                      # shared TypeScript interfaces used across features
```

---

## Detailed documentation

Read the relevant doc before implementing. These are kept in sync with the code — prefer them over general Next.js knowledge when they conflict.

| Topic | File |
|---|---|
| App Router, route files, layouts, navigation | [`docs/routing.md`](docs/routing.md) |
| Server Components, data fetching, Server Actions | [`docs/data-fetching.md`](docs/data-fetching.md) |
| Tailwind CSS v4, theme tokens, dark mode | [`docs/styling.md`](docs/styling.md) |
| Component conventions, TypeScript rules | [`docs/components.md`](docs/components.md) |
| Testing setup, patterns, what to test | [`docs/testing.md`](docs/testing.md) |

---

## Testing instructions

**TDD is required.** Write failing tests first, then implement.

### Framework
- **Vitest** + **@testing-library/react** + **jsdom** for unit and component tests.
- `next/image` and other Next.js built-ins are aliased to lightweight mocks in `vitest.config.ts`.
- Server Component rendering requires Playwright E2E (not yet set up) — test logic separately.

### What to test
- All functions in `lib/` must have Vitest unit tests.
- All Client Components (`"use client"`) must have render tests.
- Server Components: extract and test any logic; rendering is covered by E2E.

### Patterns
```tsx
// Component test
import { render, screen } from '@testing-library/react'
import MyButton from '@/components/MyButton'

it('renders label', () => {
  render(<MyButton label="Save" />)
  expect(screen.getByRole('button', { name: 'Save' })).toBeInTheDocument()
})

// Utility test
import { formatDate } from '@/lib/format'

it('formats ISO date', () => {
  expect(formatDate('2026-01-15')).toBe('Jan 15, 2026')
})
```

### Quality gate — all must pass before committing
```bash
pnpm test     # no failing tests
pnpm lint     # no ESLint errors
pnpm build    # no TypeScript errors, no build failures
```

See [`docs/testing.md`](docs/testing.md) for full setup details and conventions.

---

## Key conventions (short version — see docs for full detail)

- **App Router only** — never create a `pages/` directory
- **Server Components by default** — add `"use client"` only for browser APIs or React hooks; push it to the smallest possible component
- **Tailwind CSS v4** — no CSS modules, no inline `style={}`, no other CSS frameworks
- **No `any`** — use proper interfaces or `unknown`
- **Data fetching** — fetch in Server Components or Server Actions, not in client-side `useEffect`
- **Backend calls** — go through Server Components or Server Actions; never fetch `localhost:8080` from client components
