# Web Docs Index

Topic-based documentation for the Next.js web app. Each file covers one concern.
The `docs` agent reads this index first to locate the right file.

| Topic | File | Source files covered |
|---|---|---|
| App Router structure & route conventions | [routing.md](routing.md) | `app/layout.tsx`, `app/page.tsx`, `next.config.ts` |
| Data fetching patterns | [data-fetching.md](data-fetching.md) | `app/page.tsx`, `app/layout.tsx` |
| Styling with Tailwind CSS v4 | [styling.md](styling.md) | `app/globals.css`, `postcss.config.mjs` |
| Component conventions | [components.md](components.md) | `app/` (all component files) |
| Testing patterns | [testing.md](testing.md) | `vitest.config.ts`, `vitest.setup.ts`, `__tests__/page.test.tsx` |
| Observability (Sentry error tracking) | [observability.md](observability.md) | `sentry.client.config.ts`, `sentry.server.config.ts`, `sentry.edge.config.ts`, `next.config.ts`, `.env.example` |
| WebSocket hook (useWebSocket, reconnect, auth) | [websocket.md](websocket.md) | `lib/useWebSocket.ts`, `lib/useWebSocket.test.ts` |
| Firebase Cloud Messaging — permission, token, service worker, useFCM hook | [fcm.md](fcm.md) | `lib/fcm.ts`, `lib/useFCM.ts`, `public/firebase-messaging-sw.js`, `lib/fcm.test.ts` |
| Object storage (Cloudflare R2) — presign utility, uploadToR2, useUpload hook | [storage.md](storage.md) | `lib/storage.ts`, `lib/useUpload.ts` |
