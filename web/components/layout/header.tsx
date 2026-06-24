import Link from "next/link"
import { NavAuth } from "./NavAuth"

export default function Header() {
  return (
    <header className="sticky top-0 z-50 w-full border-b border-border bg-background/95 backdrop-blur supports-backdrop-filter:bg-background/60">
      <div className="container mx-auto flex h-14 items-center px-4">
        <Link href="/" className="text-sm font-semibold tracking-tight">
          App
        </Link>
        <div className="flex-1" />
        <NavAuth />
      </div>
    </header>
  )
}
