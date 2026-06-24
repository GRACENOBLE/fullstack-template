import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { GoogleSignInButton } from '../components/GoogleSignInButton'

vi.mock('firebase/auth', () => ({
  signInWithPopup: vi.fn(),
  GoogleAuthProvider: class {
    constructor() {}
  },
}))

vi.mock('@/lib/firebase', () => ({
  getFirebaseAuth: vi.fn(),
}))

vi.mock('next-auth/react', () => ({
  signIn: vi.fn(),
}))

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: vi.fn() }),
}))

describe('GoogleSignInButton', () => {
  it('renders the Google sign-in button', () => {
    render(<GoogleSignInButton />)
    expect(screen.getByRole('button', { name: /continue with google/i })).toBeInTheDocument()
  })
})
