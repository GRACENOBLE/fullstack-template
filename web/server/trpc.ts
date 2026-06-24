import { initTRPC, TRPCError } from '@trpc/server'
import { type NextRequest } from 'next/server'
import { auth } from '@/auth'
import type { Session } from 'next-auth'

export interface TRPCContext {
  req: NextRequest
  session: Session | null
}

export async function createTRPCContext({ req }: { req: NextRequest }): Promise<TRPCContext> {
  const session = await auth()
  return { req, session }
}

const t = initTRPC.context<TRPCContext>().create()

export const router = t.router
export const createCallerFactory = t.createCallerFactory
export const publicProcedure = t.procedure
export const protectedProcedure = t.procedure.use(async ({ ctx, next }) => {
  if (!ctx.session?.user) {
    throw new TRPCError({ code: 'UNAUTHORIZED' })
  }
  return next({ ctx })
})
