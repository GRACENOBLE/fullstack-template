---
topic: styling
last_verified: 2026-06-23
sources:
  - app/globals.css
  - postcss.config.mjs
---

# Styling

## Framework
Tailwind CSS v4. No `tailwind.config.js` — v4 uses CSS-first configuration.

## Imports
`app/globals.css` imports three stylesheets in this order:
```css
@import "tailwindcss";
@import "tw-animate-css";
@import "shadcn/tailwind.css";
```

`globals.css` must be imported in `app/layout.tsx`.

## Dark mode
Dark mode uses the `.dark` class selector, not a media query. The custom variant is declared at the top of `globals.css`:
```css
@custom-variant dark (&:is(.dark *));
```
Add the `dark` class to the `<html>` element to activate dark mode. Do not use `@media (prefers-color-scheme: dark)`.

## Theme tokens (`@theme inline`)
All design tokens are declared in `@theme inline` in `globals.css`. They map Tailwind utility names to the raw CSS custom properties defined in `:root` and `.dark`:

**Color tokens** (each has a `-foreground` counterpart where applicable):
- `--color-background`, `--color-foreground`
- `--color-card`, `--color-card-foreground`
- `--color-popover`, `--color-popover-foreground`
- `--color-primary`, `--color-primary-foreground`
- `--color-secondary`, `--color-secondary-foreground`
- `--color-muted`, `--color-muted-foreground`
- `--color-accent`, `--color-accent-foreground`
- `--color-destructive`
- `--color-border`, `--color-input`, `--color-ring`
- `--color-chart-1` through `--color-chart-5`
- `--color-sidebar`, `--color-sidebar-foreground`, `--color-sidebar-primary`, `--color-sidebar-primary-foreground`, `--color-sidebar-accent`, `--color-sidebar-accent-foreground`, `--color-sidebar-border`, `--color-sidebar-ring`

**Radius tokens** (derived from base `--radius: 0.625rem`):
- `--radius-sm` = `calc(var(--radius) * 0.6)`
- `--radius-md` = `calc(var(--radius) * 0.8)`
- `--radius-lg` = `var(--radius)`
- `--radius-xl` = `calc(var(--radius) * 1.4)`
- `--radius-2xl` = `calc(var(--radius) * 1.8)`
- `--radius-3xl` = `calc(var(--radius) * 2.2)`
- `--radius-4xl` = `calc(var(--radius) * 2.6)`

**Font tokens:**
- `--font-sans: var(--font-sans)`
- `--font-mono: var(--font-geist-mono)`
- `--font-heading: var(--font-sans)`

## Color space
All raw color values in `:root` and `.dark` use `oklch(...)`. Example:
```css
:root {
  --primary: oklch(0.205 0 0);
  --destructive: oklch(0.577 0.245 27.325);
}
```

## Base layer
`@layer base` in `globals.css` applies defaults globally:
```css
@layer base {
  * {
    @apply border-border outline-ring/50;
  }
  body {
    @apply bg-background text-foreground;
  }
  html {
    @apply font-sans;
  }
}
```

## Toaster theming
The `.toaster` block wires shadcn/sonner toast colors to the shared token system:
```css
.toaster {
  --normal-bg: var(--popover);
  --normal-text: var(--popover-foreground);
  --normal-border: var(--border);
  --border-radius: var(--radius);
}
```

## PostCSS
Config in `postcss.config.mjs`. Uses `@tailwindcss/postcss` plugin. Do not modify unless adding a non-Tailwind PostCSS plugin.

## Rules
- Tailwind classes only — no CSS modules, no styled-components, no inline `style={}`.
- New design tokens go in `globals.css` under `@theme inline`, not in a config file.
- Always use semantic tokens (`bg-background`, `text-foreground`, `bg-primary`) over raw color utilities (`bg-white`, `text-gray-900`) so dark mode works automatically.
- Do not use `@media (prefers-color-scheme: dark)` — the project uses class-based dark mode exclusively.
