'use client'

import { signOut as nextAuthSignOut } from 'next-auth/react'
import { useRouter } from 'next/navigation'

export function useSignOut() {
  const router = useRouter()

  const signOut = async () => {
    await nextAuthSignOut({ redirect: false })
    router.push('/login')
  }

  return { signOut }
}
