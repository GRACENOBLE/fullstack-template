import Link from 'next/link'

export function Logo({ className }: { className?: string }) {
  return (
    <Link href="/" className={`text-sm font-semibold tracking-tight hover:opacity-70 transition-opacity ${className ?? ''}`}>
      App
    </Link>
  )
}
