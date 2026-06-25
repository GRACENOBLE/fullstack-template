"use client"

import { Button } from "@/components/ui/button"

interface ErrorPageProps {
  error: Error & { digest?: string }
  reset: () => void
}

export default function ErrorPage({ error, reset }: ErrorPageProps) {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center gap-6 px-4 text-center">
      <div className="space-y-2">
        <h1 className="text-3xl font-bold text-foreground">
          Something went wrong
        </h1>
        <p className="text-muted-foreground">
          {error.message || "An unexpected error occurred. Please try again."}
        </p>
        {error.digest && (
          <p
            data-testid="error-digest"
            className="text-xs text-muted-foreground font-mono"
          >
            Reference: {error.digest}
          </p>
        )}
      </div>
      <Button onClick={reset}>Try again</Button>
    </div>
  )
}
