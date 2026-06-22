---
topic: storage
last_verified: 2026-06-23
sources:
  - web/lib/storage.ts
  - web/lib/useUpload.ts
---

# Storage (Cloudflare R2 upload)

## Upload flow
The client never talks to R2 directly with credentials. The two-step flow is:

1. Call the backend `POST /api/v1/storage/presign` with the filename and MIME type to receive a short-lived presigned PUT URL and the final public URL.
2. PUT the file bytes directly to R2 using the presigned URL.

Both steps are encapsulated in `lib/storage.ts`. The `useUpload` hook in `lib/useUpload.ts` wraps the flow with React state.

## lib/storage.ts

### `presign(filename, contentType, idToken): Promise<PresignResult>`
Posts to the backend presign endpoint. Forwards the Firebase ID token as `Authorization: Bearer <idToken>`. Throws if the response is not OK.

Return type:
```ts
interface PresignResult {
  uploadUrl: string
  publicUrl: string
}
```

The backend base URL is read from `process.env.NEXT_PUBLIC_BACKEND_URL`, defaulting to `http://localhost:8080`.

Request body sent:
```json
{ "filename": "avatar.png", "content_type": "image/png" }
```

### `uploadToR2(file: File, uploadUrl: string): Promise<void>`
PUTs the file directly to R2 using the presigned URL. Sets `Content-Type` to `file.type`. Throws if the response is not OK. No auth header — the presigned URL is self-authenticating.

## lib/useUpload.ts

`useUpload` is a Client Component hook (`'use client'`). It orchestrates the two-step upload and exposes React state.

```ts
interface UploadState {
  isUploading: boolean
  publicUrl: string | null
  error: Error | null
  upload: (file: File, idToken: string) => Promise<void>
}

function useUpload(): UploadState
```

Calling `upload(file, idToken)`:
1. Sets `isUploading = true`, clears `error`.
2. Calls `presign(file.name, file.type, idToken)` to get `uploadUrl` and `publicUrl`.
3. Calls `uploadToR2(file, uploadUrl)`.
4. On success: sets `publicUrl`.
5. On any error: sets `error` (normalised to `Error` if the thrown value is not already one).
6. Always sets `isUploading = false` in the `finally` block.

The `upload` function is memoised with `useCallback` and has no dependencies, so it is stable across renders.

## Environment variable

| Variable | Description |
|---|---|
| `NEXT_PUBLIC_BACKEND_URL` | Backend base URL. Defaults to `http://localhost:8080` when absent. |

## Testing

### storage.ts
Tests use `vi.stubGlobal('fetch', ...)` to intercept the `presign` call and the R2 PUT without any real network. Assert on the URL, method, headers, and body passed to `fetch`.

### useUpload.ts
Tests use `vi.mock('@/lib/storage', ...)` to replace `presign` and `uploadToR2` with controlled stubs. Render the hook with `renderHook`, call `result.current.upload(...)`, and assert on `isUploading`, `publicUrl`, and `error`.
