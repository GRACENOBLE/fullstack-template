import { initTRPC, TRPCError } from '@trpc/server'
import { type NextRequest } from 'next/server'

export interface TRPCContext {
  req: NextRequest
  // session will be added when auth is implemented
}

export async function createTRPCContext({ req }: { req: NextRequest }): Promise<TRPCContext> {
  return { req }
}

const t = initTRPC.context<TRPCContext>().create()

export const router = t.router
export const createCallerFactory = t.createCallerFactory
export const publicProcedure = t.procedure
export const protectedProcedure = t.procedure.use(async ({ ctx, next }) => {
  const authHeader = ctx.req.headers.get('authorization')
  const cookieHeader = ctx.req.headers.get('cookie')

  const bearerToken = authHeader?.startsWith('Bearer ') ? authHeader.slice(7).trim() : null
  const hasValidBearer = typeof bearerToken === 'string' && bearerToken.length > 0

  const sessionValue = cookieHeader
    ?.split(';')
    .map(c => c.trim())
    .find(c => c.startsWith('__session='))
    ?.slice('__session='.length)
    .trim()
  const hasValidSessionCookie = typeof sessionValue === 'string' && sessionValue.length > 0

  if (!hasValidBearer && !hasValidSessionCookie) {
    throw new TRPCError({ code: 'UNAUTHORIZED' })
  }

  return next({ ctx })
})
