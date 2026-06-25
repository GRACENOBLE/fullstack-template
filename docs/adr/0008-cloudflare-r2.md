# ADR 0008 — Cloudflare R2 for object storage

**Status:** Accepted  
**Date:** 2026-06-25

## Context

The template needs file/object storage for user uploads (avatars, attachments, exports). Requirements:
- S3-compatible API to avoid vendor-specific client code
- No egress fees for download traffic
- Opt-in — projects that don't need file storage should not pay for it
- Presigned URL support so the backend does not proxy file bytes

Candidates evaluated: Cloudflare R2, AWS S3, Supabase Storage.

## Decision

Use **Cloudflare R2** accessed via the **AWS SDK v2** (`github.com/aws/aws-sdk-go-v2/service/s3`), configured with R2's S3-compatible endpoint (`https://<account-id>.r2.cloudflarestorage.com`).

The storage feature is **opt-in**: omitting `R2_ACCOUNT_ID`, `R2_ACCESS_KEY`, and `R2_SECRET_KEY` from the environment disables the `/api/v1/storage/presign` and `/api/v1/storage/:key` routes. The `Handler` struct's `storageService` field will be `nil`, and `RegisterRoutes` conditionally skips registration.

Upload flow:
1. Client calls `POST /api/v1/storage/presign` → backend returns a presigned S3 PUT URL (valid 15 min).
2. Client uploads directly to R2 using the presigned URL — backend never proxies bytes.
3. The public URL (`R2_PUBLIC_URL`) is returned for download access.

## Consequences

### Positive
- Zero egress fees from R2 regardless of download volume — significant cost saving at scale compared to S3.
- The AWS SDK v2 means switching to S3 or any S3-compatible provider (MinIO, Backblaze B2) requires only environment variable changes, no code changes.
- Presigned URLs eliminate the backend as a bottleneck for large file uploads.
- The opt-in pattern means projects without file storage have zero R2 surface area.

### Negative / trade-offs
- Cloudflare account required. The free R2 tier covers 10 GB storage and 1M Class A operations/month.
- `R2_PUBLIC_URL` must be configured separately (custom domain or `*.r2.dev` subdomain); the bucket is private by default.
- No local emulator for R2 is provided — integration tests for the storage layer require either a live R2 bucket or MinIO running locally.
