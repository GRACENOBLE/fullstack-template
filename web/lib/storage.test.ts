import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'

import { presign, uploadToR2 } from './storage'

beforeEach(() => {
  vi.clearAllMocks()
})

afterEach(() => {
  vi.unstubAllGlobals()
})

describe('presign', () => {
  it('returns uploadUrl and publicUrl on a 200 response', async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        upload_url: 'https://r2.example.com/presigned-put',
        public_url: 'https://cdn.example.com/my-file.jpg',
      }),
    })
    vi.stubGlobal('fetch', mockFetch)

    const result = await presign('my-file.jpg', 'image/jpeg', 'id-token-abc')

    expect(result).toEqual({
      uploadUrl: 'https://r2.example.com/presigned-put',
      publicUrl: 'https://cdn.example.com/my-file.jpg',
    })
  })

  it('throws when the backend returns a non-200 response', async () => {
    const mockFetch = vi.fn().mockResolvedValue({ ok: false, status: 401 })
    vi.stubGlobal('fetch', mockFetch)

    await expect(presign('file.png', 'image/png', 'bad-token')).rejects.toThrow('presign failed: 401')
  })

  it('passes the Authorization header with the idToken', async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ upload_url: 'https://r2.example.com/put', public_url: 'https://cdn.example.com/f' }),
    })
    vi.stubGlobal('fetch', mockFetch)

    await presign('photo.png', 'image/png', 'my-id-token')

    expect(mockFetch).toHaveBeenCalledWith(
      expect.stringContaining('/api/v1/storage/presign'),
      expect.objectContaining({
        headers: expect.objectContaining({
          Authorization: 'Bearer my-id-token',
        }),
      }),
    )
  })

  it('sends filename and content_type in the request body', async () => {
    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ upload_url: 'https://r2.example.com/put', public_url: 'https://cdn.example.com/f' }),
    })
    vi.stubGlobal('fetch', mockFetch)

    await presign('doc.pdf', 'application/pdf', 'token-xyz')

    expect(mockFetch).toHaveBeenCalledWith(
      expect.anything(),
      expect.objectContaining({
        method: 'POST',
        body: JSON.stringify({ filename: 'doc.pdf', content_type: 'application/pdf' }),
      }),
    )
  })
})

describe('uploadToR2', () => {
  it('makes a PUT request to the presigned URL with the file bytes and Content-Type', async () => {
    const mockFetch = vi.fn().mockResolvedValue({ ok: true })
    vi.stubGlobal('fetch', mockFetch)

    const file = new File(['hello world'], 'hello.txt', { type: 'text/plain' })
    await uploadToR2(file, 'https://r2.example.com/presigned-put?sig=abc')

    expect(mockFetch).toHaveBeenCalledWith(
      'https://r2.example.com/presigned-put?sig=abc',
      expect.objectContaining({
        method: 'PUT',
        headers: expect.objectContaining({ 'Content-Type': 'text/plain' }),
        body: file,
      }),
    )
  })

  it('throws on a non-2xx response from R2', async () => {
    const mockFetch = vi.fn().mockResolvedValue({ ok: false, status: 403 })
    vi.stubGlobal('fetch', mockFetch)

    const file = new File(['data'], 'data.bin', { type: 'application/octet-stream' })
    await expect(uploadToR2(file, 'https://r2.example.com/put')).rejects.toThrow('R2 upload failed: 403')
  })
})
