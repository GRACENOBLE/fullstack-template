import { z } from 'zod'
import { protectedProcedure, router } from '../trpc'

export const notificationsRouter = router({
  registerFcmToken: protectedProcedure
    .input(z.object({ token: z.string().min(1) }))
    .mutation(async ({ input }) => {
      // Stub — will be replaced when notifications are implemented
      return { registered: true, token: input.token }
    }),
  list: protectedProcedure.query(async () => {
    // Stub — will be replaced when notifications are implemented
    return []
  }),
})
