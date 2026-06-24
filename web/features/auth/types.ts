export interface AuthUser {
  id?: string | null
  name?: string | null
  email?: string | null
  image?: string | null
}

export interface AuthSession {
  user: AuthUser
  expires: string
}
