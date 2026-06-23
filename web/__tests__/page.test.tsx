import { render } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import Home from '../app/page'

describe('Home page', () => {
  it('renders without crashing', () => {
    const { container } = render(<Home />)
    expect(container.firstChild).toBeInTheDocument()
  })
})
