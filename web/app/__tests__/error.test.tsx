import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import ErrorPage from '../error'

describe('ErrorPage', () => {
  it('renders a friendly error heading', () => {
    const error = new Error('Something went wrong')
    const reset = vi.fn()
    render(<ErrorPage error={error} reset={reset} />)
    expect(screen.getByRole('heading')).toBeInTheDocument()
  })

  it('renders the error message in a paragraph', () => {
    const error = new Error('Network connection lost')
    const reset = vi.fn()
    render(<ErrorPage error={error} reset={reset} />)
    expect(screen.getByText(/network connection lost/i)).toBeInTheDocument()
  })

  it('calls reset when the Try again button is clicked', async () => {
    const user = userEvent.setup()
    const error = new Error('Boom')
    const reset = vi.fn()
    render(<ErrorPage error={error} reset={reset} />)
    await user.click(screen.getByRole('button', { name: /try again/i }))
    expect(reset).toHaveBeenCalledTimes(1)
  })

  it('shows the digest when present', () => {
    const error = Object.assign(new Error('Boom'), { digest: 'abc-123' })
    const reset = vi.fn()
    render(<ErrorPage error={error} reset={reset} />)
    expect(screen.getByText(/abc-123/)).toBeInTheDocument()
  })

  it('does not render digest text when digest is absent', () => {
    const error = new Error('Boom')
    const reset = vi.fn()
    render(<ErrorPage error={error} reset={reset} />)
    expect(screen.queryByTestId('error-digest')).not.toBeInTheDocument()
  })
})
