import { cert, getApp, getApps, initializeApp } from 'firebase-admin/app'
import { getAuth } from 'firebase-admin/auth'

function getAdminApp() {
  if (getApps().length > 0) return getApp()

  const serviceAccountJson = process.env.FIREBASE_SERVICE_ACCOUNT_JSON
  if (serviceAccountJson) {
    return initializeApp({
      credential: cert(JSON.parse(serviceAccountJson) as object),
    })
  }

  const projectId = process.env.FIREBASE_PROJECT_ID
  if (!projectId) {
    throw new Error(
      'Firebase Admin: set FIREBASE_SERVICE_ACCOUNT_JSON or FIREBASE_PROJECT_ID',
    )
  }
  return initializeApp({ projectId })
}

export async function verifyFirebaseToken(idToken: string) {
  return getAuth(getAdminApp()).verifyIdToken(idToken)
}
