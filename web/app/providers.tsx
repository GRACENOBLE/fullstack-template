'use client'

import { TRPCProvider } from '@/lib/trpc/client'

export function Providers({ children }: { children: React.ReactNode }) {
  return <TRPCProvider>{children}</TRPCProvider>
}
