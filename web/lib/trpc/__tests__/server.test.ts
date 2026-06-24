import { describe, it, expect, vi } from 'vitest'

vi.mock('next/headers', () => ({
  headers: vi.fn(async () => new Headers({ 'x-test': 'true' })),
}))

vi.mock('react', async (importOriginal) => {
  const actual = await importOriginal<typeof import('react')>()
  return { ...actual, cache: (fn: unknown) => fn }
})

const { createServerCaller } = await import('../server')

describe('createServerCaller', () => {
  it('returns a caller with the expected router shape', async () => {
    const caller = await createServerCaller()
    expect(typeof caller.health.query).toBe('function')
    expect(typeof caller.auth.session).toBe('function')
    expect(typeof caller.auth.signOut).toBe('function')
    expect(typeof caller.notifications.list).toBe('function')
  })

  it('propagates request headers into the tRPC context', async () => {
    const { headers } = await import('next/headers')
    const mockHeaders = new Headers({ authorization: 'Bearer test-token' })
    vi.mocked(headers).mockResolvedValueOnce(mockHeaders as Awaited<ReturnType<typeof headers>>)
    const caller = await createServerCaller()
    expect(caller).toBeDefined()
  })
})
