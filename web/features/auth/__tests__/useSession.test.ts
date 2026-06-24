import { describe, it, expect, vi } from 'vitest'
import { renderHook } from '@testing-library/react'
import { useSession } from '../hooks/useSession'

vi.mock('next-auth/react', () => ({
  useSession: vi.fn(),
}))

import { useSession as useNextAuthSession } from 'next-auth/react'
const mockUseSession = vi.mocked(useNextAuthSession)

describe('useSession', () => {
  it('returns unauthenticated state when no session', () => {
    mockUseSession.mockReturnValue({
      data: null,
      status: 'unauthenticated',
      update: vi.fn(),
    })

    const { result } = renderHook(() => useSession())

    expect(result.current.isAuthenticated).toBe(false)
    expect(result.current.isLoading).toBe(false)
    expect(result.current.user).toBeNull()
  })

  it('returns loading state', () => {
    mockUseSession.mockReturnValue({
      data: null,
      status: 'loading',
      update: vi.fn(),
    })

    const { result } = renderHook(() => useSession())

    expect(result.current.isLoading).toBe(true)
    expect(result.current.isAuthenticated).toBe(false)
  })

  it('returns authenticated state with user', () => {
    mockUseSession.mockReturnValue({
      data: {
        user: { id: '1', name: 'Grace Noble', email: 'grace@example.com', image: null },
        expires: '2099-01-01',
      },
      status: 'authenticated',
      update: vi.fn(),
    })

    const { result } = renderHook(() => useSession())

    expect(result.current.isAuthenticated).toBe(true)
    expect(result.current.isLoading).toBe(false)
    expect(result.current.user?.email).toBe('grace@example.com')
    expect(result.current.user?.name).toBe('Grace Noble')
  })
})
