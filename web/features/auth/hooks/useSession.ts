'use client'

import { useSession as useNextAuthSession } from 'next-auth/react'
import type { AuthSession, AuthUser } from '../types'

export interface UseSessionReturn {
  user: AuthUser | null
  isAuthenticated: boolean
  isLoading: boolean
}

export function useSession(): UseSessionReturn {
  const { data, status } = useNextAuthSession()

  return {
    user: (data as AuthSession | null)?.user ?? null,
    isAuthenticated: status === 'authenticated',
    isLoading: status === 'loading',
  }
}
