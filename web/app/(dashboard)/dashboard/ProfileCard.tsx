import type { UserProfile } from '@/lib/user-profile'

interface ProfileCardProps {
  profile: UserProfile
}

export function ProfileCard({ profile }: ProfileCardProps) {
  return (
    <div className="rounded-lg border bg-card p-6 shadow-sm">
      <div className="space-y-1">
        <p className="text-2xl font-bold text-foreground">{profile.displayName}</p>
        <p className="text-base text-muted-foreground">{profile.email}</p>
        <p className="font-mono text-xs text-muted-foreground/70">{profile.uid}</p>
      </div>
    </div>
  )
}
