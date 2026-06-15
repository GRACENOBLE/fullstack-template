importScripts('https://www.gstatic.com/firebasejs/11.0.0/firebase-app-compat.js')
importScripts('https://www.gstatic.com/firebasejs/11.0.0/firebase-messaging-compat.js')

// Firebase config is injected by the page before registering this worker,
// via postMessage, or you can hardcode NEXT_PUBLIC_* values here at build time.
// For development, populate these from your .env.local file.
self.addEventListener('message', (event) => {
  if (event.data?.type === 'FIREBASE_CONFIG') {
    firebase.initializeApp(event.data.config)
    firebase.messaging()
  }
})

// Background message handler — called when the app is not in the foreground.
self.addEventListener('push', () => {
  // firebase-messaging-compat handles push events automatically once initialised.
})
