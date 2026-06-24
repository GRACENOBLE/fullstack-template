import { describe, it, expect, vi, beforeEach } from 'vitest'
import { appRouter } from '../_app'
import { createTRPCContext } from '../../trpc'

// mock fetch
const mockFetch = vi.fn()
vi.stubGlobal('fetch', mockFetch)

function makeContext(reqHeaders: Record<string, string> = {}) {
  const req = new Request('http://localhost/api/trpc', { headers: reqHeaders }) as import('next/server').NextRequest
  return createTRPCContext({ req })
}

describe('health router', () => {
  beforeEach(() => {
    mockFetch.mockReset()
  })

  it('returns health data from backend', async () => {
    mockFetch.mockResolvedValueOnce(
      new Response(JSON.stringify({ status: 'ok', database: 'ok' }), { status: 200 })
    )
    const ctx = await makeContext()
    const caller = appRouter.createCaller(ctx)
    const result = await caller.health.query()
    expect(result).toEqual({ status: 'ok', database: 'ok' })
    expect(mockFetch).toHaveBeenCalledWith(expect.stringContaining('/health'))
  })

  it('throws when backend is down', async () => {
    mockFetch.mockResolvedValueOnce(new Response('', { status: 503 }))
    const ctx = await makeContext()
    const caller = appRouter.createCaller(ctx)
    await expect(caller.health.query()).rejects.toThrow('Backend health check failed: 503')
  })
})
