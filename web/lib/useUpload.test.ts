import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, act } from '@testing-library/react'

const { mockPresign, mockUploadToR2 } = vi.hoisted(() => ({
  mockPresign: vi.fn(),
  mockUploadToR2: vi.fn(),
}))

vi.mock('@/lib/storage', () => ({
  presign: mockPresign,
  uploadToR2: mockUploadToR2,
}))

import { useUpload } from './useUpload'

beforeEach(() => {
  vi.clearAllMocks()
})

describe('useUpload', () => {
  it('has isUploading false initially', () => {
    const { result } = renderHook(() => useUpload())
    expect(result.current.isUploading).toBe(false)
  })

  it('has publicUrl null initially', () => {
    const { result } = renderHook(() => useUpload())
    expect(result.current.publicUrl).toBeNull()
  })

  it('has error null initially', () => {
    const { result } = renderHook(() => useUpload())
    expect(result.current.error).toBeNull()
  })

  it('sets publicUrl and clears isUploading after a successful upload', async () => {
    mockPresign.mockResolvedValue({
      uploadUrl: 'https://r2.example.com/put',
      publicUrl: 'https://cdn.example.com/uploaded.jpg',
    })
    mockUploadToR2.mockResolvedValue(undefined)

    const { result } = renderHook(() => useUpload())
    const file = new File(['img'], 'photo.jpg', { type: 'image/jpeg' })

    await act(async () => {
      await result.current.upload(file, 'id-token-ok')
    })

    expect(result.current.publicUrl).toBe('https://cdn.example.com/uploaded.jpg')
    expect(result.current.isUploading).toBe(false)
    expect(result.current.error).toBeNull()
  })

  it('sets error and clears isUploading after a failed upload', async () => {
    mockPresign.mockRejectedValue(new Error('presign failed: 401'))

    const { result } = renderHook(() => useUpload())
    const file = new File(['img'], 'photo.jpg', { type: 'image/jpeg' })

    await act(async () => {
      await result.current.upload(file, 'bad-token')
    })

    expect(result.current.error).toBeInstanceOf(Error)
    expect(result.current.error?.message).toBe('presign failed: 401')
    expect(result.current.isUploading).toBe(false)
    expect(result.current.publicUrl).toBeNull()
  })

  it('sets isUploading to true while the upload is in progress', async () => {
    let resolvePresign!: (value: { uploadUrl: string; publicUrl: string }) => void
    mockPresign.mockReturnValue(
      new Promise<{ uploadUrl: string; publicUrl: string }>((resolve) => {
        resolvePresign = resolve
      }),
    )
    mockUploadToR2.mockResolvedValue(undefined)

    const { result } = renderHook(() => useUpload())
    const file = new File(['img'], 'photo.jpg', { type: 'image/jpeg' })

    // Start the upload without awaiting so we can inspect mid-flight state
    act(() => {
      void result.current.upload(file, 'token')
    })

    expect(result.current.isUploading).toBe(true)

    // Resolve and finish
    await act(async () => {
      resolvePresign({ uploadUrl: 'https://r2.example.com/put', publicUrl: 'https://cdn.example.com/f' })
    })

    expect(result.current.isUploading).toBe(false)
  })

  it('calls presign with the file name, file type, and idToken', async () => {
    mockPresign.mockResolvedValue({
      uploadUrl: 'https://r2.example.com/put',
      publicUrl: 'https://cdn.example.com/f',
    })
    mockUploadToR2.mockResolvedValue(undefined)

    const { result } = renderHook(() => useUpload())
    const file = new File(['data'], 'report.pdf', { type: 'application/pdf' })

    await act(async () => {
      await result.current.upload(file, 'auth-token')
    })

    expect(mockPresign).toHaveBeenCalledWith('report.pdf', 'application/pdf', 'auth-token')
  })

  it('calls uploadToR2 with the file and the presigned upload URL', async () => {
    mockPresign.mockResolvedValue({
      uploadUrl: 'https://r2.example.com/put?sig=xyz',
      publicUrl: 'https://cdn.example.com/f',
    })
    mockUploadToR2.mockResolvedValue(undefined)

    const { result } = renderHook(() => useUpload())
    const file = new File(['data'], 'image.png', { type: 'image/png' })

    await act(async () => {
      await result.current.upload(file, 'token')
    })

    expect(mockUploadToR2).toHaveBeenCalledWith(file, 'https://r2.example.com/put?sig=xyz')
  })
})
