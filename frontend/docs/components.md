---
topic: components
last_verified: 2026-06-14
sources:
  - app/layout.tsx
  - app/page.tsx
---

# Component Conventions

## Default: Server Components
All components are Server Components unless they have a `"use client"` directive.
Server Components can be `async`, fetch data directly, and access server-only resources.

## When to add `"use client"`
Only when you need:
- `useState`, `useReducer`, `useContext`, or other React hooks
- `useEffect` or lifecycle behavior
- Browser APIs (`window`, `document`, `localStorage`)
- Event handlers (`onClick`, `onChange`, etc.) on the component itself

Do not add `"use client"` to layouts, pages, or wrapper components just because a child needs it — push `"use client"` down to the smallest possible component.

## File placement
```
app/               # route files only (page.tsx, layout.tsx, loading.tsx, error.tsx)
components/        # shared reusable components
  ui/              # generic primitives (Button, Input, Modal, etc.)
  [feature]/       # feature-specific components (e.g., components/auth/)
lib/               # utility functions, helpers, non-React code
types/             # shared TypeScript type definitions
```

## Component file structure
```tsx
// components/ui/Button.tsx
import type { ButtonHTMLAttributes } from 'react';

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: 'primary' | 'secondary';
};

export function Button({ variant = 'primary', children, ...props }: ButtonProps) {
  return (
    <button
      className={variant === 'primary' ? 'bg-foreground text-background' : 'bg-background text-foreground border'}
      {...props}
    >
      {children}
    </button>
  );
}
```

## TypeScript rules
- No `any`. Use proper interfaces or `unknown`.
- Extend native HTML element types (`ButtonHTMLAttributes`, `InputHTMLAttributes`) for wrapper components.
- Keep component prop types in the same file as the component unless shared across multiple files.
- Shared types → `types/` directory.

## Imports
- Use TypeScript path aliases from `tsconfig.json` — check it before writing relative import chains.
- Import React only when needed (React 19 — JSX transform is automatic, no `import React` needed).
- Import types with `import type` to keep runtime bundles clean.
