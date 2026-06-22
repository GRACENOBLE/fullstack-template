import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { renderHook } from '@testing-library/react'

const { mockRequestPermissionAndGetToken, mockOnForegroundMessage } = vi.hoisted(() => ({
  mockRequestPermissionAndGetToken: vi.fn(),
  mockOnForegroundMessage: vi.fn(() => vi.fn()),
}))

vi.mock('@/lib/fcm', () => ({
  requestPermissionAndGetToken: mockRequestPermissionAndGetToken,
  onForegroundMessage: mockOnForegroundMessage,
}))

import { useFCM } from './useFCM'

beforeEach(() => {
  vi.clearAllMocks()
  mockOnForegroundMessage.mockReturnValue(vi.fn())
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('useFCM', () => {
  it('does not register when token is null', async () => {
    mockRequestPermissionAndGetToken.mockResolvedValue(null)
    const fetchSpy = vi.fn()
    vi.stubGlobal('fetch', fetchSpy)

    const { unmount } = renderHook(() => useFCM())
    // Allow the promise microtask to settle
    await Promise.resolve()

    expect(fetchSpy).not.toHaveBeenCalled()
    unmount()
  })

  it('registers token when permission granted', async () => {
    mockRequestPermissionAndGetToken.mockResolvedValue('test-token')
    const fetchSpy = vi.fn().mockResolvedValue({ ok: true })
    vi.stubGlobal('fetch', fetchSpy)

    const { unmount } = renderHook(() => useFCM())
    await Promise.resolve()

    expect(fetchSpy).toHaveBeenCalledWith(
      '/api/v1/fcm/register',
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ token: 'test-token', platform: 'web' }),
      }),
    )
    unmount()
  })

  it('passes Authorization header when idToken provided', async () => {
    mockRequestPermissionAndGetToken.mockResolvedValue('test-token')
    const fetchSpy = vi.fn().mockResolvedValue({ ok: true })
    vi.stubGlobal('fetch', fetchSpy)

    const { unmount } = renderHook(() => useFCM({ idToken: 'id-tok' }))
    await Promise.resolve()

    expect(fetchSpy).toHaveBeenCalledWith(
      '/api/v1/fcm/register',
      expect.objectContaining({
        headers: expect.objectContaining({
          Authorization: 'Bearer id-tok',
        }),
      }),
    )
    unmount()
  })

  it('subscribes to foreground messages when onMessage provided', async () => {
    mockRequestPermissionAndGetToken.mockResolvedValue(null)
    const handler = vi.fn()

    const { unmount } = renderHook(() => useFCM({ onMessage: handler }))
    await Promise.resolve()

    expect(mockOnForegroundMessage).toHaveBeenCalledWith(handler)
    unmount()
  })

  it('unsubscribes on unmount', async () => {
    mockRequestPermissionAndGetToken.mockResolvedValue(null)
    const unsub = vi.fn()
    mockOnForegroundMessage.mockReturnValue(unsub)
    const handler = vi.fn()

    const { unmount } = renderHook(() => useFCM({ onMessage: handler }))
    await Promise.resolve()

    unmount()
    expect(unsub).toHaveBeenCalled()
  })
})
