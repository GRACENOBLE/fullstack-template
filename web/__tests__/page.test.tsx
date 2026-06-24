import { render } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import Home from '../app/page'

vi.mock('next-auth/react', () => ({
  useSession: vi.fn(() => ({ data: null, status: 'unauthenticated' })),
  SessionProvider: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}))

describe('Home page', () => {
  it('renders without crashing', () => {
    const { container } = render(<Home />)
    expect(container.firstChild).toBeInTheDocument()
  })
})
