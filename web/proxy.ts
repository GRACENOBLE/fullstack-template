import { auth } from "@/auth"
import { NextResponse } from "next/server"
import type { NextAuthRequest } from "next-auth"

export default auth((req: NextAuthRequest) => {
  if (!req.auth) {
    const loginUrl = new URL("/login", req.url)
    loginUrl.searchParams.set("callbackUrl", req.nextUrl.pathname + req.nextUrl.search)
    return NextResponse.redirect(loginUrl)
  }
})

export const config = {
  matcher: ["/dashboard/:path*", "/settings/:path*"],
}
