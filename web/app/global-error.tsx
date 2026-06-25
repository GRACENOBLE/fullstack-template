"use client"

import { Button } from "@/components/ui/button"

interface GlobalErrorProps {
  error: Error & { digest?: string }
  reset: () => void
}

export default function GlobalError({ error, reset }: GlobalErrorProps) {
  return (
    <html lang="en" className="h-full antialiased">
      <body className="min-h-full flex flex-col bg-background text-foreground font-sans">
        <div className="min-h-screen flex flex-col items-center justify-center gap-6 px-4 text-center">
          <div className="space-y-2">
            <h1 className="text-3xl font-bold text-foreground">
              Something went wrong
            </h1>
            <p className="text-muted-foreground">
              {error.message ||
                "A critical error occurred. Please reload the page."}
            </p>
            {error.digest && (
              <p className="text-xs text-muted-foreground font-mono">
                Reference: {error.digest}
              </p>
            )}
          </div>
          <Button onClick={reset}>Try again</Button>
        </div>
      </body>
    </html>
  )
}
