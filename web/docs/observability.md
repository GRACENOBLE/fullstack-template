---
topic: observability
last_verified: 2026-06-15
sources:
  - sentry.client.config.ts
  - sentry.server.config.ts
  - sentry.edge.config.ts
  - next.config.ts
  - .env.example
---

# Observability

## Sentry SDK

| Package | Version |
|---|---|
| `@sentry/nextjs` | 10.57.0 |

## Runtime config files

Three files initialize Sentry, each targeting a different Next.js runtime. All three call `Sentry.init` with identical options:

```ts
Sentry.init({
  dsn: process.env.NEXT_PUBLIC_SENTRY_DSN,
  tracesSampleRate: 1.0,
  debug: false,
});
```

| File | Runtime |
|---|---|
| `sentry.client.config.ts` | Browser (client components, client-side navigation) |
| `sentry.server.config.ts` | Node.js (Server Components, Route Handlers, Server Actions) |
| `sentry.edge.config.ts` | Edge runtime (Middleware, Edge Route Handlers) |

Next.js picks up each file automatically based on the execution context — no manual imports are needed.

## Environment variable

| Variable | Required | Default |
|---|---|---|
| `NEXT_PUBLIC_SENTRY_DSN` | No | `undefined` (Sentry disabled) |

The `NEXT_PUBLIC_` prefix makes the value available in the browser bundle. When the variable is absent or empty, `Sentry.init` is a no-op.

**Local development** — copy `web/.env.example` to `web/.env.local` and fill in your DSN:
```dotenv
NEXT_PUBLIC_SENTRY_DSN=https://<key>@o<org>.ingest.sentry.io/<project>
```

**Production** — set `NEXT_PUBLIC_SENTRY_DSN` in your deployment platform's environment configuration.

## next.config.ts — withSentryConfig

`next.config.ts` wraps the Next.js config with `withSentryConfig`:

```ts
export default withSentryConfig(nextConfig, {
  silent: true,
  org: "",      // fill in your Sentry org slug
  project: "",  // fill in your Sentry project slug
  widenClientFileUpload: true,
  sourcemaps: {
    deleteSourcemapsAfterUpload: true,
  },
  webpack: {
    treeshake: {
      removeDebugLogging: true,
    },
  },
});
```

Option notes:
- `silent: true` — suppresses Sentry CLI output during builds.
- `org` and `project` — must be set to non-empty strings to enable source map uploads to Sentry. Left as `""` in the committed file; fill these in for production builds.
- `widenClientFileUpload: true` — uploads source maps from a broader set of client bundle chunks.
- `sourcemaps.deleteSourcemapsAfterUpload: true` — removes source map files from the build output after uploading, so they are not served publicly.
- `webpack.treeshake.removeDebugLogging: true` — strips Sentry debug log calls from the production bundle.
