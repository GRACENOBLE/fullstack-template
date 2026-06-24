"use client"

import Link from "next/link"
import { useSession } from "@/features/auth/hooks/useSession"
import { UserMenu } from "@/features/auth/components/UserMenu"
import { Button } from "@/components/ui/button"

export function NavAuth() {
  const { isAuthenticated, isLoading } = useSession()

  if (isLoading) return <div className="h-8 w-8 rounded-full bg-muted animate-pulse" />
  if (isAuthenticated) return <UserMenu />

  return (
    <Button asChild variant="ghost" size="sm">
      <Link href="/login">Sign in</Link>
    </Button>
  )
}
