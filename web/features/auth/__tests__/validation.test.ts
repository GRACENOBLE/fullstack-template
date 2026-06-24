import { describe, it, expect } from 'vitest'
import { loginSchema, registerSchema } from '../validation'

describe('loginSchema', () => {
  it('validates correct credentials', () => {
    const result = loginSchema.safeParse({ email: 'user@example.com', password: 'password123' })
    expect(result.success).toBe(true)
  })

  it('rejects invalid email', () => {
    const result = loginSchema.safeParse({ email: 'not-an-email', password: 'password123' })
    expect(result.success).toBe(false)
    expect(result.error?.issues[0]?.path).toContain('email')
  })

  it('rejects empty password', () => {
    const result = loginSchema.safeParse({ email: 'user@example.com', password: '' })
    expect(result.success).toBe(false)
    expect(result.error?.issues[0]?.path).toContain('password')
  })

  it('rejects missing fields', () => {
    const result = loginSchema.safeParse({})
    expect(result.success).toBe(false)
  })
})

describe('registerSchema', () => {
  const valid = {
    name: 'Grace Noble',
    email: 'grace@example.com',
    password: 'securepassword',
    confirmPassword: 'securepassword',
  }

  it('validates correct registration data', () => {
    expect(registerSchema.safeParse(valid).success).toBe(true)
  })

  it('rejects name shorter than 2 characters', () => {
    const result = registerSchema.safeParse({ ...valid, name: 'G' })
    expect(result.success).toBe(false)
    expect(result.error?.issues[0]?.path).toContain('name')
  })

  it('rejects password shorter than 8 characters', () => {
    const result = registerSchema.safeParse({
      ...valid,
      password: 'short',
      confirmPassword: 'short',
    })
    expect(result.success).toBe(false)
    expect(result.error?.issues[0]?.path).toContain('password')
  })

  it('rejects mismatched passwords', () => {
    const result = registerSchema.safeParse({ ...valid, confirmPassword: 'different' })
    expect(result.success).toBe(false)
    expect(result.error?.issues[0]?.path).toContain('confirmPassword')
  })

  it('rejects invalid email', () => {
    const result = registerSchema.safeParse({ ...valid, email: 'not-email' })
    expect(result.success).toBe(false)
  })
})
