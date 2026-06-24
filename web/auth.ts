import NextAuth from "next-auth"
import Credentials from "next-auth/providers/credentials"
import { z } from "zod"

export const { handlers, auth, signIn, signOut } = NextAuth({
  providers: [
    Credentials({
      credentials: { idToken: {} },
      async authorize(credentials) {
        const parsed = z.object({ idToken: z.string().min(1) }).safeParse(credentials)
        if (!parsed.success) return null

        const parts = parsed.data.idToken.split(".")
        if (parts.length !== 3) return null

        try {
          const payload = JSON.parse(
            Buffer.from(parts[1], "base64url").toString("utf-8"),
          ) as { sub?: string; email?: string; name?: string; picture?: string }

          if (!payload.sub) return null

          return {
            id: payload.sub,
            email: payload.email ?? null,
            name: payload.name ?? null,
            image: payload.picture ?? null,
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
