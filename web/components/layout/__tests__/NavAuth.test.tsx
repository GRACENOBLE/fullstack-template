import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { NavAuth } from '../NavAuth'

vi.mock('@/features/auth/hooks/useSession', () => ({
  useSession: vi.fn(),
}))

vi.mock('@/features/auth/components/UserMenu', () => ({
  UserMenu: () => <div data-testid="user-menu" />,
}))

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: vi.fn() }),
}))

import { useSession } from '@/features/auth/hooks/useSession'
const mockUseSession = vi.mocked(useSession)

describe('NavAuth', () => {
  it('shows a loading skeleton while session is loading', () => {
    mockUseSession.mockReturnValue({ isLoading: true, isAuthenticated: false, session: null })
    const { container } = render(<NavAuth />)
    expect(container.querySelector('.animate-pulse')).toBeInTheDocument()
  })

  it('renders UserMenu when authenticated', () => {
    mockUseSession.mockReturnValue({ isLoading: false, isAuthenticated: true, session: { user: { email: 'a@b.com' }, expires: '2099' } })
    render(<NavAuth />)
    expect(screen.getByTestId('user-menu')).toBeInTheDocument()
  })

  it('renders a sign-in link when unauthenticated', () => {
    mockUseSession.mockReturnValue({ isLoading: false, isAuthenticated: false, session: null })
    render(<NavAuth />)
    expect(screen.getByRole('link', { name: /sign in/i })).toBeInTheDocument()
  })
})
