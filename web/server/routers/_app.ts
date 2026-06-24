import { router } from '../trpc'
import { authRouter } from './auth'
import { healthRouter } from './health'
import { notificationsRouter } from './notifications'

export const appRouter = router({
  health: healthRouter,
  auth: authRouter,
  notifications: notificationsRouter,
})

export type AppRouter = typeof appRouter
