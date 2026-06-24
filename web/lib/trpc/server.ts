import { createCallerFactory } from '@/server/trpc'
import { cache } from 'react'
import { headers } from 'next/headers'
import { createTRPCContext } from '@/server/trpc'
import { appRouter } from '@/server/routers/_app'

const createCaller = createCallerFactory(appRouter)

export const createServerCaller = cache(async () => {
  const headerList = await headers()
  // Build a minimal Request-like object for the context
  const req = new Request('http://internal', {
    headers: headerList,
  }) as import('next/server').NextRequest
  const ctx = await createTRPCContext({ req })
  return createCaller(ctx)
})
