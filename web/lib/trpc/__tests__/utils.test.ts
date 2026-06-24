// @vitest-environment node
import { describe, it, expect, afterEach } from 'vitest'
import { getBaseUrl } from '../utils'

afterEach(() => {
  delete process.env.VERCEL_URL
})

describe('getBaseUrl', () => {
  it('returns empty string when window is defined (browser)', () => {
    // In jsdom/browser environments window is defined; here we simulate it in node
    const g = global as Record<string, unknown>
    g.window = {}
    try {
      expect(getBaseUrl()).toBe('')
    } finally {
      delete g.window
    }
  })

  it('returns Vercel URL when VERCEL_URL env var is set', () => {
    process.env.VERCEL_URL = 'my-app.vercel.app'
    expect(getBaseUrl()).toBe('https://my-app.vercel.app')
  })

  it('falls back to localhost:3000', () => {
    expect(getBaseUrl()).toBe('http://localhost:3000')
  })
})
