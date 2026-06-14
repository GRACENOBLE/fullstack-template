<!-- BEGIN:nextjs-agent-rules -->
# This is NOT the Next.js you know

This version has breaking changes — APIs, conventions, and file structure may all differ from your training data. Read the relevant guide in `node_modules/next/dist/docs/` before writing any code. Heed deprecation notices.
<!-- END:nextjs-agent-rules -->

---

## Commands

```bash
pnpm install    # install dependencies
pnpm dev        # dev server → http://localhost:3000
pnpm build      # production build (also runs TypeScript check)
pnpm start      # serve production build
pnpm lint       # ESLint
```

---

## Project structure

```
app/                  # Next.js App Router — routes only
  layout.tsx          # root layout: Geist fonts, global CSS, <html>/<body>
  page.tsx            # home page (Server Component)
  globals.css         # Tailwind v4 import + CSS variable definitions
public/               # static assets
components/           # shared UI components (create when needed)
lib/                  # utilities and non-React helpers (create when needed)
types/                # shared TypeScript type definitions (create when needed)
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

---

## Testing instructions

No component test suite yet. Quality gate is lint + build:

```bash
pnpm lint     # must pass — no ESLint errors
pnpm build    # must pass — no TypeScript errors, no build failures
```

Both must pass before committing or opening a PR.

---

## Key conventions (short version — see docs for full detail)

- **App Router only** — never create a `pages/` directory
- **Server Components by default** — add `"use client"` only for browser APIs or React hooks; push it to the smallest possible component
- **Tailwind CSS v4** — no CSS modules, no inline `style={}`, no other CSS frameworks
- **No `any`** — use proper interfaces or `unknown`
- **Data fetching** — fetch in Server Components or Server Actions, not in client-side `useEffect`
- **Backend calls** — go through Server Components or Server Actions; never fetch `localhost:8080` from client components
