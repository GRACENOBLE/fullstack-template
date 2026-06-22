'use client'

import { useState, useCallback } from 'react'
import { presign, uploadToR2 } from '@/lib/storage'

export interface UploadState {
  isUploading: boolean
  publicUrl: string | null
  error: Error | null
  upload: (file: File, idToken: string) => Promise<void>
}

/**
 * React hook for the presign → PUT → public URL upload flow.
 *
 * Call `upload(file, idToken)` to start an upload. The hook tracks
 * isUploading, publicUrl (set on success), and error (set on failure).
 */
export function useUpload(): UploadState {
  const [isUploading, setIsUploading] = useState(false)
  const [publicUrl, setPublicUrl] = useState<string | null>(null)
  const [error, setError] = useState<Error | null>(null)

  const upload = useCallback(async (file: File, idToken: string) => {
    setIsUploading(true)
    setError(null)
    try {
      const { uploadUrl, publicUrl: url } = await presign(file.name, file.type, idToken)
      await uploadToR2(file, uploadUrl)
      setPublicUrl(url)
    } catch (err) {
      setError(err instanceof Error ? err : new Error(String(err)))
    } finally {
      setIsUploading(false)
    }
  }, [])

  return { isUploading, publicUrl, error, upload }
}
