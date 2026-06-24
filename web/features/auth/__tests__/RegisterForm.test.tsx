import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { RegisterForm } from '../components/RegisterForm'

vi.mock('firebase/auth', () => ({
  createUserWithEmailAndPassword: vi.fn(),
  updateProfile: vi.fn(),
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

describe('RegisterForm', () => {
  it('renders all registration fields', () => {
    render(<RegisterForm />)
    expect(screen.getByLabelText(/^name$/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/^email$/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/^password$/i)).toBeInTheDocument()
    expect(screen.getByLabelText(/confirm password/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /create account/i })).toBeInTheDocument()
  })

  it('shows validation error when passwords do not match', async () => {
    render(<RegisterForm />)
    const user = userEvent.setup()

    await user.type(screen.getByLabelText(/^name$/i), 'Grace Noble')
    await user.type(screen.getByLabelText(/^email$/i), 'grace@example.com')
    await user.type(screen.getByLabelText(/^password$/i), 'password123')
    await user.type(screen.getByLabelText(/confirm password/i), 'differentpassword')
    await user.click(screen.getByRole('button', { name: /create account/i }))

    expect(await screen.findByText(/passwords don't match/i)).toBeInTheDocument()
  })

  it('shows validation error for short password', async () => {
    render(<RegisterForm />)
    const user = userEvent.setup()

    await user.type(screen.getByLabelText(/^name$/i), 'Grace Noble')
    await user.type(screen.getByLabelText(/^email$/i), 'grace@example.com')
    await user.type(screen.getByLabelText(/^password$/i), 'short')
    await user.type(screen.getByLabelText(/confirm password/i), 'short')
    await user.click(screen.getByRole('button', { name: /create account/i }))

    expect(await screen.findByText(/at least 8 characters/i)).toBeInTheDocument()
  })
})
