import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { fetchUserProfile } from './user-profile'

const originalEnv = process.env

beforeEach(() => {
  process.env = { ...originalEnv, BACKEND_URL: 'http://localhost:8080' }
})

afterEach(() => {
  process.env = originalEnv
  vi.restoreAllMocks()
})

describe('fetchUserProfile', () => {
  it('returns a UserProfile when the backend responds with 200', async () => {
    const fakeProfile = {
      uid: 'firebase-uid-123',
      email: 'alice@example.com',
      displayName: 'Alice Example',
    }

    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValueOnce({
        ok: true,
        json: async () => ({ data: fakeProfile }),
      }),
    )

    const result = await fetchUserProfile('firebase-uid-123')
    expect(result).toEqual(fakeProfile)
  })

  it('sends the X-User-Id header with the provided userId', async () => {
    const mockFetch = vi.fn().mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        data: { uid: 'u1', email: 'u@e.com', displayName: 'U' },
      }),
    })
    vi.stubGlobal('fetch', mockFetch)

    await fetchUserProfile('my-user-id')

    expect(mockFetch).toHaveBeenCalledWith(
      'http://localhost:8080/api/v1/me',
      expect.objectContaining({
        headers: expect.objectContaining({ 'X-User-Id': 'my-user-id' }),
      }),
    )
  })

  it('throws when the backend returns a non-2xx status', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValueOnce({
        ok: false,
        status: 401,
        statusText: 'Unauthorized',
      }),
    )

    await expect(fetchUserProfile('uid')).rejects.toThrow(
      'Failed to fetch user profile: 401 Unauthorized',
    )
  })

  it('throws when BACKEND_URL is not set', async () => {
    delete process.env.BACKEND_URL

    await expect(fetchUserProfile('uid')).rejects.toThrow(
      'BACKEND_URL environment variable is not set',
    )
  })

  it('throws when the network call fails', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockRejectedValueOnce(new Error('Network error')),
    )

    await expect(fetchUserProfile('uid')).rejects.toThrow('Network error')
  })
})
