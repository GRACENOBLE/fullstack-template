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
  // Read session token from Authorization header (Bearer) or __session cookie
  const authHeader = ctx.req.headers.get('authorization')
  const cookieHeader = ctx.req.headers.get('cookie')

  const hasBearer = authHeader?.startsWith('Bearer ')
  const hasSessionCookie = cookieHeader?.includes('__session=')

  if (!hasBearer && !hasSessionCookie) {
    throw new TRPCError({ code: 'UNAUTHORIZED' })
  }

  return next({ ctx })
})
