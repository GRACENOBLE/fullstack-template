import { auth } from '@/auth'
import { redirect } from 'next/navigation'
import { fetchUserProfile } from '@/lib/user-profile'
import type { UserProfile } from '@/lib/user-profile'
import { ProfileCard } from './ProfileCard'

export default async function DashboardPage() {
  const session = await auth()
  if (!session) redirect('/login')

  // Build a fallback profile from session data in case the backend fetch fails.
  // Note: NextAuth is configured with the JWT strategy using a Credentials provider
  // that verifies the Firebase ID token server-side during sign-in. The original
  // Firebase ID token is not forwarded into the NextAuth session — only the decoded
  // claims are stored (id, email, name). The backend fetch therefore cannot use a
  // Firebase Bearer token; it sends X-User-Id instead.
  const fallbackProfile: UserProfile = {
    uid: session.user?.id ?? '',
    email: session.user?.email ?? '',
    displayName: session.user?.name ?? session.user?.email ?? 'User',
  }

  let profile: UserProfile = fallbackProfile
  if (fallbackProfile.uid) {
    try {
      profile = await fetchUserProfile(fallbackProfile.uid)
    } catch {
      // Backend unreachable or returned a non-2xx status — show session data instead.
      profile = fallbackProfile
    }
  }

  return (
    <div className="flex flex-1 flex-col gap-6 p-8">
      <div>
        <h1 className="text-3xl font-bold">Welcome back</h1>
        <p className="mt-1 text-muted-foreground">Here is your profile information.</p>
      </div>
      <ProfileCard profile={profile} />
    </div>
  )
}
