import NextAuth from "next-auth"
import Credentials from "next-auth/providers/credentials"
import { z } from "zod"
import { verifyFirebaseToken } from "@/lib/firebase-admin"

export const { handlers, auth, signIn, signOut } = NextAuth({
  providers: [
    Credentials({
      credentials: { idToken: {} },
      async authorize(credentials) {
        const parsed = z.object({ idToken: z.string().min(1) }).safeParse(credentials)
        if (!parsed.success) return null

        try {
          const decoded = await verifyFirebaseToken(parsed.data.idToken)
          if (!decoded.sub) return null

          return {
            id: decoded.sub,
            email: decoded.email ?? null,
            name: decoded.name ?? null,
            image: decoded.picture ?? null,
          }
        } catch {
          return null
        }
      },
    }),
  ],
  pages: { signIn: "/login" },
  session: { strategy: "jwt" },
})
