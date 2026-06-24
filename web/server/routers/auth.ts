import { protectedProcedure, router } from '../trpc'

export const authRouter = router({
  session: protectedProcedure.query(async ({ ctx }) => {
    return {
      authenticated: true,
      user: ctx.session?.user ?? null,
    }
  }),
  signOut: protectedProcedure.mutation(async () => {
    return { success: true }
  }),
})
