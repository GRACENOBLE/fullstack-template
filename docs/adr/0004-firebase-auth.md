# ADR 0004 — Firebase Authentication (cross-platform)

**Status:** Accepted  
**Date:** 2026-06-25

## Context

The template needs an authentication system that:
- Works across all three layers: Go backend, Next.js web, and Android
- Supports Google OAuth and email/password out of the box
- Handles token lifecycle (refresh, revocation) without server-side session storage
- Has a generous free tier suitable for early-stage projects

Candidates evaluated: Firebase Auth, Auth0, Supabase Auth, custom JWT.

## Decision

Use **Firebase Authentication** for identity management across all three layers.

**Flow:**
1. User signs in via the Firebase client SDK (web or Android).
2. Firebase issues a short-lived ID token (1 hour).
3. The client sends the ID token as `Authorization: Bearer <token>` to the Go backend.
4. The Go backend verifies the token using the Firebase Admin SDK (`firebase.google.com/go/v4`).
5. On the web, NextAuth wraps Firebase sign-in with a server-side session; only the decoded claims (uid, email, name) are stored in the NextAuth JWT — the raw Firebase ID token is not forwarded.

**Error responses** from the backend FirebaseAuth middleware use the standard envelope:
```json
{"error": {"code": "UNAUTHORIZED", "message": "..."}}
```

## Consequences

### Positive
- Google OAuth, email/password, Apple, and phone auth are all available with no custom code.
- Firebase ID tokens are JWTs verifiable offline once the public keys are cached — no Firebase round-trip per request.
- Cross-platform SDK consistency: the same Firebase project serves web and Android.
- The Firebase Console provides user management, analytics, and remote config at no cost at typical template scale.

### Negative / trade-offs
- Vendor lock-in to Google Firebase. Migrating to a self-hosted solution requires replacing the Admin SDK verification and all client-side sign-in flows.
- The NextAuth/Firebase bridge means the backend cannot receive a Bearer token from the web layer — the web sends `X-User-Id` as an identifier instead of a verifiable token. A future improvement would store a Firebase custom token in the NextAuth session.
- `google-services.json` (Android) and `GoogleService-Info.plist` (iOS, if added) are project-specific secrets and must not be committed.
