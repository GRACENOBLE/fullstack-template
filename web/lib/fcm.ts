'use client'

import { initializeApp, getApps, type FirebaseApp } from 'firebase/app'
import { getMessaging, getToken, onMessage, type MessagePayload } from 'firebase/messaging'

const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
  storageBucket: process.env.NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET,
  messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID,
}

function getFirebaseApp(): FirebaseApp {
  if (getApps().length > 0) return getApps()[0]!
  return initializeApp(firebaseConfig)
}

/**
 * Requests Notification permission, registers the service worker, and returns
 * the FCM registration token. Returns null when permission is denied or the
 * browser does not support notifications.
 */
export async function requestPermissionAndGetToken(): Promise<string | null> {
  if (typeof window === 'undefined' || !('Notification' in window)) return null

  const permission = await Notification.requestPermission()
  if (permission !== 'granted') return null

  const registration = await navigator.serviceWorker.register('/firebase-messaging-sw.js')
  const messaging = getMessaging(getFirebaseApp())

  return getToken(messaging, {
    vapidKey: process.env.NEXT_PUBLIC_FIREBASE_VAPID_KEY,
    serviceWorkerRegistration: registration,
  })
}

/**
 * Subscribes to foreground FCM messages. Returns an unsubscribe function.
 * Call this once per component lifecycle; call the returned function on cleanup.
 */
export function onForegroundMessage(handler: (payload: MessagePayload) => void): () => void {
  const messaging = getMessaging(getFirebaseApp())
  return onMessage(messaging, handler)
}
