---
topic: fcm
last_verified: 2026-06-16
sources:
  - lib/fcm.ts
  - lib/useFCM.ts
  - public/firebase-messaging-sw.js
  - lib/fcm.test.ts
---

# Firebase Cloud Messaging (FCM) — Web

Enables browser push notifications via the Firebase JS SDK. Works in foreground (tab open) and background (tab closed, via service worker).

## Environment variables

Add these to `.env.local` (never commit real values):

```
NEXT_PUBLIC_FIREBASE_API_KEY=
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=
NEXT_PUBLIC_FIREBASE_PROJECT_ID=
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=
NEXT_PUBLIC_FIREBASE_APP_ID=
NEXT_PUBLIC_FIREBASE_VAPID_KEY=
```

Obtain the web app config from **Firebase Console → Project Settings → General → Your apps**. Obtain the VAPID key from **Firebase Console → Project Settings → Cloud Messaging → Web Push certificates**.

## lib/fcm.ts

Pure utility module. Must only be imported in Client Components (`"use client"`).

```ts
// Request notification permission and obtain the FCM registration token.
// Returns null if the browser does not support notifications or if the
// user denies the permission prompt.
requestPermissionAndGetToken(): Promise<string | null>

// Subscribe to foreground messages. Returns an unsubscribe function.
// Call the returned function in the component's cleanup.
onForegroundMessage(handler: (payload: MessagePayload) => void): () => void
```

`requestPermissionAndGetToken` registers `public/firebase-messaging-sw.js` as a service worker before calling `getToken`, so the service worker is ready before the FCM token is obtained.

Firebase app initialisation is lazy and idempotent: `getApps().length > 0` check prevents re-initialising on hot reloads.

## lib/useFCM.ts

Client hook. Requests permission, gets the FCM token, and POSTs it to the backend on mount:

```ts
useFCM({ idToken?: string; onMessage?: (payload: MessagePayload) => void }): void
```

- `idToken`: Firebase Auth ID token. Sent as `Authorization: Bearer <token>` to authenticate the backend registration request. Omit when Firebase Auth is not yet wired on the web.
- `onMessage`: called for each foreground message while the component is mounted.

The hook re-runs (re-registers the token) whenever `idToken` changes, so switching users automatically re-registers the token under the new identity.

Example usage in a root layout or auth-protected page:
```tsx
'use client'
import { useFCM } from '@/lib/useFCM'

export default function NotificationProvider({ idToken }: { idToken?: string }) {
  useFCM({
    idToken,
    onMessage: (payload) => console.log('FCM foreground:', payload),
  })
  return null
}
```

## public/firebase-messaging-sw.js

Service worker for background messages (tab closed or not focused). Loaded by `requestPermissionAndGetToken` automatically.

The service worker receives the Firebase config from the page via `postMessage` with `{ type: 'FIREBASE_CONFIG', config: { ... } }`. Send this immediately after registering the worker:

```ts
const reg = await navigator.serviceWorker.register('/firebase-messaging-sw.js')
reg.active?.postMessage({ type: 'FIREBASE_CONFIG', config: firebaseConfig })
```

## Testing

Tests live in `lib/fcm.test.ts` (Vitest, jsdom). Both `firebase/app` and `firebase/messaging` are fully mocked with `vi.mock` + `vi.hoisted()`.

Covered paths:
- Returns `null` when `Notification` is absent from `window`
- Returns `null` when the user denies permission
- Returns the FCM token on success (verifies `getToken` is called with the service worker registration)
- `onForegroundMessage` registers the handler and returns the unsubscribe function
