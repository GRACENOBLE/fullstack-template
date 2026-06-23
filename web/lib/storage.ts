const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL ?? 'http://localhost:8080'

export interface PresignResult {
  uploadUrl: string
  publicUrl: string
}

/**
 * Calls the backend presign endpoint to get a presigned PUT URL and a public URL.
 * The idToken is forwarded as a Bearer token so the backend can authenticate the caller.
 */
export async function presign(
  filename: string,
  contentType: string,
  idToken: string,
): Promise<PresignResult> {
  const res = await fetch(`${BACKEND_URL}/api/v1/storage/presign`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${idToken}`,
    },
    body: JSON.stringify({ filename, content_type: contentType }),
  })
  if (!res.ok) {
    throw new Error(`presign failed: ${res.status}`)
  }
  const data = (await res.json()) as { upload_url: string; public_url: string }
  return { uploadUrl: data.upload_url, publicUrl: data.public_url }
}

/**
 * Uploads a file directly to Cloudflare R2 using a presigned PUT URL.
 * The file bytes are sent as the request body with the file's MIME type.
 */
export async function uploadToR2(file: File, uploadUrl: string): Promise<void> {
  const res = await fetch(uploadUrl, {
    method: 'PUT',
    headers: { 'Content-Type': file.type },
    body: file,
  })
  if (!res.ok) {
    throw new Error(`R2 upload failed: ${res.status}`)
  }
}
