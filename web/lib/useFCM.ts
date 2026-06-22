'use client'

import { useEffect } from 'react'
import { requestPermissionAndGetToken, onForegroundMessage } from '@/lib/fcm'
import type { MessagePayload } from 'firebase/messaging'

/**
 * Requests push notification permission, obtains the FCM token, and registers
 * it with the backend. Optionally subscribes to foreground messages.
 *
 * idToken: Firebase Auth ID token to authenticate the registration request.
 * onMessage: called for each foreground push message received.
 */
export function useFCM({
  idToken,
  onMessage,
}: {
  idToken?: string
  onMessage?: (payload: MessagePayload) => void
} = {}): void {
  useEffect(() => {
    let cancelled = false

    requestPermissionAndGetToken().then(async (token) => {
      if (!token || cancelled) return
      try {
        const response = await fetch('/api/v1/fcm/register', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            ...(idToken ? { Authorization: `Bearer ${idToken}` } : {}),
          },
          body: JSON.stringify({ token, platform: 'web' }),
        })
        if (!response.ok) {
          console.error('FCM token registration failed:', response.status)
        }
      } catch (err) {
        console.error('FCM token registration error:', err)
      }
    })

    const unsub = onMessage ? onForegroundMessage(onMessage) : null
    return () => {
      cancelled = true
      unsub?.()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [idToken])
}
