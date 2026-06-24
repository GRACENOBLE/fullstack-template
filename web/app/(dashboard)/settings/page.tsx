import { auth } from '@/auth'
import { redirect } from 'next/navigation'

export default async function SettingsPage() {
  const session = await auth()
  if (!session) redirect('/login')

  return (
    <div className="flex flex-1 flex-col gap-6 p-8">
      <div>
        <h1 className="text-2xl font-bold">Settings</h1>
        <p className="mt-1 text-muted-foreground">Manage your account settings.</p>
      </div>
      <div className="rounded-lg border p-6">
        <h2 className="text-sm font-medium">Account</h2>
        <p className="mt-1 text-sm text-muted-foreground">
          {session.user?.name && (
            <span className="block">{session.user.name}</span>
          )}
          {session.user?.email}
        </p>
      </div>
    </div>
  )
}
