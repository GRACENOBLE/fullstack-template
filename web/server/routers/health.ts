import { publicProcedure, router } from '../trpc'

const BACKEND_URL = process.env.BACKEND_URL ?? 'http://localhost:8080'

export const healthRouter = router({
  query: publicProcedure.query(async () => {
    const res = await fetch(`${BACKEND_URL}/health`)
    if (!res.ok) throw new Error(`Backend health check failed: ${res.status}`)
    return res.json() as Promise<{ status: string; database: string }>
  }),
})
