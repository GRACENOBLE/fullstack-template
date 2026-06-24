import { protectedProcedure, router } from '../trpc'

export const authRouter = router({
  session: protectedProcedure.query(async () => {
    // Stub — will be replaced when auth is implemented
    return { authenticated: true }
  }),
  signOut: protectedProcedure.mutation(async () => {
    // Stub — will be replaced when auth is implemented
    return { success: true }
  }),
})
