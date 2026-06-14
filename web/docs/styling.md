---
topic: styling
last_verified: 2026-06-14
sources:
  - app/globals.css
  - postcss.config.mjs
---

# Styling

## Framework
Tailwind CSS v4. No `tailwind.config.js` — v4 uses CSS-first configuration.

## Import
```css
/* app/globals.css — must be imported in app/layout.tsx */
@import "tailwindcss";
```

## Theme customization (CSS variables via `@theme inline`)
Semantic design tokens are defined in `globals.css` using `@theme inline`:
```css
@theme inline {
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --font-sans: var(--font-geist-sans);
  --font-mono: var(--font-geist-mono);
}
```
These tokens become Tailwind utilities: `bg-background`, `text-foreground`, `font-sans`, `font-mono`.

## Dark mode
CSS variable swap via media query in `globals.css`:
```css
@media (prefers-color-scheme: dark) {
  :root {
    --background: #0a0a0a;
    --foreground: #ededed;
  }
}
```
Use `bg-background` and `text-foreground` Tailwind classes — they automatically follow dark mode.

## Fonts
Geist Sans and Geist Mono loaded in `layout.tsx` via `next/font/google`.
They inject `--font-geist-sans` and `--font-geist-mono` CSS variables.
These are wired into Tailwind's `font-sans` and `font-mono` utilities via `@theme inline`.

## Rules
- Tailwind classes only — no CSS modules, no styled-components, no inline `style={}`.
- Add new design tokens in `globals.css` under `@theme inline`, not in a config file.
- Use semantic tokens (`bg-background`, `text-foreground`) over raw colors (`bg-white`, `text-gray-900`) so dark mode works automatically.
- For component-specific styles that can't be handled by Tailwind, add to `globals.css` with a clear comment. Keep it minimal.

## PostCSS
Config in `postcss.config.mjs`. Uses `@tailwindcss/postcss` plugin. Do not modify unless adding a non-Tailwind PostCSS plugin.
