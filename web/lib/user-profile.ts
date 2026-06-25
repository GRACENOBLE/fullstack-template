/**
 * Fetches the authenticated user's profile from the Go backend /api/v1/me endpoint.
 *
 * Auth note: NextAuth is configured with the JWT strategy using a Credentials provider
 * that verifies the Firebase ID token server-side during sign-in. The original Firebase
 * ID token is NOT forwarded into the NextAuth session/JWT — only the decoded claims
 * (uid, email, name) are stored. As a result, this function cannot send a Firebase
 * Bearer token to the backend. It uses the NextAuth user ID as a fallback identifier
 * in the X-User-Id header. A future improvement would be to store a long-lived custom
 * token or re-issue a Firebase custom token in a NextAuth JWT callback.
 */

export interface UserProfile {
  uid: string
  email: string
  displayName: string
}

interface ApiEnvelope {
  data: UserProfile
}

export async function fetchUserProfile(userId: string): Promise<UserProfile> {
  const backendUrl = process.env.BACKEND_URL
  if (!backendUrl) {
    throw new Error('BACKEND_URL environment variable is not set')
  }

  const res = await fetch(`${backendUrl}/api/v1/me`, {
    cache: 'no-store',
    headers: {
      'X-User-Id': userId,
    },
  })

  if (!res.ok) {
    throw new Error(`Failed to fetch user profile: ${res.status} ${res.statusText}`)
  }

  const body: ApiEnvelope = (await res.json()) as ApiEnvelope
  return body.data
}
