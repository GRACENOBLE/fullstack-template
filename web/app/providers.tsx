'use client'

import { TRPCProvider } from '@/lib/trpc/client'
import { SessionProvider } from 'next-auth/react'
import { Toaster } from '@/components/ui/sonner'

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <SessionProvider>
      <TRPCProvider>{children}</TRPCProvider>
      <Toaster richColors position="top-right" />
    </SessionProvider>
  )
}
