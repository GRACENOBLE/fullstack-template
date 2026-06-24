'use client'

import { signOut as nextAuthSignOut } from 'next-auth/react'
import { signOut as firebaseSignOut } from 'firebase/auth'
import { useRouter } from 'next/navigation'
import { getFirebaseAuth } from '@/lib/firebase'

export function useSignOut() {
  const router = useRouter()

  const signOut = async () => {
    await Promise.all([
      nextAuthSignOut({ redirect: false }),
      firebaseSignOut(getFirebaseAuth()),
    ])
    router.push('/login')
  }

  return { signOut }
}
