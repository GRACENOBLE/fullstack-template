import { describe, it, expect, vi, beforeEach } from 'vitest'

const { mockGetToken, mockOnMessage, mockGetMessaging, mockGetApps, mockInitializeApp } =
  vi.hoisted(() => ({
    mockGetToken: vi.fn(),
    mockOnMessage: vi.fn(() => vi.fn()),
    mockGetMessaging: vi.fn(() => ({})),
    mockGetApps: vi.fn(() => [] as unknown[]),
    mockInitializeApp: vi.fn(() => ({ name: '[DEFAULT]' })),
  }))

vi.mock('firebase/app', () => ({
  initializeApp: mockInitializeApp,
  getApps: mockGetApps,
}))

vi.mock('firebase/messaging', () => ({
  getMessaging: mockGetMessaging,
  getToken: mockGetToken,
  onMessage: mockOnMessage,
}))

import { requestPermissionAndGetToken, onForegroundMessage } from './fcm'

beforeEach(() => {
  vi.clearAllMocks()
  mockGetApps.mockReturnValue([])
  mockGetMessaging.mockReturnValue({})
})

describe('requestPermissionAndGetToken', () => {
  it('returns null when Notification is not in window', async () => {
    const saved = (globalThis as Record<string, unknown>).Notification
    delete (globalThis as Record<string, unknown>).Notification

    const result = await requestPermissionAndGetToken()
    expect(result).toBeNull()
    expect(mockGetToken).not.toHaveBeenCalled()

    ;(globalThis as Record<string, unknown>).Notification = saved
  })

  it('returns null when permission is denied', async () => {
    Object.defineProperty(globalThis, 'Notification', {
      value: { requestPermission: vi.fn().mockResolvedValue('denied') },
      configurable: true,
      writable: true,
    })

    const result = await requestPermissionAndGetToken()
    expect(result).toBeNull()
    expect(mockGetToken).not.toHaveBeenCalled()
  })

  it('returns the FCM token when permission is granted', async () => {
    const fakeToken = 'fcm-registration-token-abc'
    const fakeRegistration = { scope: '/firebase-messaging-sw.js' }

    Object.defineProperty(globalThis, 'Notification', {
      value: { requestPermission: vi.fn().mockResolvedValue('granted') },
      configurable: true,
      writable: true,
    })
    Object.defineProperty(globalThis, 'navigator', {
      value: { serviceWorker: { register: vi.fn().mockResolvedValue(fakeRegistration) } },
      configurable: true,
      writable: true,
    })
    mockGetToken.mockResolvedValue(fakeToken)

    const result = await requestPermissionAndGetToken()

    expect(result).toBe(fakeToken)
    expect(mockGetToken).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({ serviceWorkerRegistration: fakeRegistration }),
    )
  })
})

describe('onForegroundMessage', () => {
  it('registers the handler and returns the unsubscribe function', () => {
    const handler = vi.fn()
    const unsubscribe = vi.fn()
    mockOnMessage.mockReturnValue(unsubscribe)

    const result = onForegroundMessage(handler)

    expect(mockOnMessage).toHaveBeenCalledWith(expect.anything(), handler)
    expect(result).toBe(unsubscribe)
  })
})
