'use client'

import { TRPCProvider } from '@/lib/trpc/client'
import { SessionProvider } from 'next-auth/react'
import { NuqsAdapter } from 'nuqs/adapters/next/app'
import { Toaster } from '@/components/ui/sonner'

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <NuqsAdapter>
      <SessionProvider>
        <TRPCProvider>{children}</TRPCProvider>
        <Toaster richColors position="top-right" />
      </SessionProvider>
    </NuqsAdapter>
  )
}
