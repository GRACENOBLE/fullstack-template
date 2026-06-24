import { describe, it, expect } from 'vitest'
import { appRouter } from '../_app'
import { createTRPCContext } from '../../trpc'

function makeContext(reqHeaders: Record<string, string> = {}) {
  const req = new Request('http://localhost/api/trpc', { headers: reqHeaders }) as import('next/server').NextRequest
  return createTRPCContext({ req })
}

describe('protectedProcedure', () => {
  it('throws UNAUTHORIZED when no session present', async () => {
    const ctx = await makeContext()
    const caller = appRouter.createCaller(ctx)
    await expect(caller.auth.session()).rejects.toMatchObject({
      code: 'UNAUTHORIZED',
    })
  })

  it('allows access with Bearer token', async () => {
    const ctx = await makeContext({ authorization: 'Bearer test-token-123' })
    const caller = appRouter.createCaller(ctx)
    const result = await caller.auth.session()
    expect(result).toEqual({ authenticated: true })
  })

  it('allows access with __session cookie', async () => {
    const ctx = await makeContext({ cookie: '__session=abc123' })
    const caller = appRouter.createCaller(ctx)
    const result = await caller.auth.session()
    expect(result).toEqual({ authenticated: true })
  })

  it('throws UNAUTHORIZED for empty Bearer token', async () => {
    const ctx = await makeContext({ authorization: 'Bearer ' })
    const caller = appRouter.createCaller(ctx)
    await expect(caller.auth.session()).rejects.toMatchObject({ code: 'UNAUTHORIZED' })
  })

  it('throws UNAUTHORIZED for __session cookie with empty value', async () => {
    const ctx = await makeContext({ cookie: '__session=' })
    const caller = appRouter.createCaller(ctx)
    await expect(caller.auth.session()).rejects.toMatchObject({ code: 'UNAUTHORIZED' })
  })

  it('throws UNAUTHORIZED when __session= appears only in another cookie value', async () => {
    const ctx = await makeContext({ cookie: 'other=__session=abc' })
    const caller = appRouter.createCaller(ctx)
    await expect(caller.auth.session()).rejects.toMatchObject({ code: 'UNAUTHORIZED' })
  })
})
