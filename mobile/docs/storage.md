---
topic: storage
last_verified: 2026-06-23
sources:
  - mobile/app/src/main/java/com/company/template/storage/UploadRepository.kt
---

# Storage (Cloudflare R2 upload)

## Upload flow
Same two-step flow as the web client: presign then PUT.

1. POST to `<backendBaseUrl>/api/v1/storage/presign` with `filename` and `content_type` to receive `upload_url` and `public_url`.
2. PUT the raw file bytes directly to R2 using the presigned URL.

Both steps execute on `Dispatchers.IO` inside a `withContext` block.

## Interface

```kotlin
interface UploadRepository {
    suspend fun upload(
        filename: String,
        contentType: String,
        fileBytes: ByteArray,
        idToken: String,
    ): Result<String>  // returns public URL on success
}
```

## R2UploadRepository

`R2UploadRepository` in `com.company.template.storage` is the concrete implementation.

Constructor:
```kotlin
class R2UploadRepository(
    private val backendBaseUrl: String,
    private val httpClient: OkHttpClient = OkHttpClient(),
)
```

The `httpClient` parameter defaults to a plain `OkHttpClient()`. Tests can inject a custom instance (e.g. one backed by `MockWebServer`) without subclassing.

### Internal data classes
Both are `@Serializable` and use `kotlinx.serialization`:

```kotlin
private data class PresignRequest(
    val filename: String,
    @SerialName("content_type") val contentType: String,
)

private data class PresignResponse(
    @SerialName("upload_url") val uploadUrl: String,
    @SerialName("public_url") val publicUrl: String,
)
```

`Json { ignoreUnknownKeys = true }` is used for decoding so extra fields from the backend do not cause failures.

### upload()
Calls private `presign()` then private `uploadToR2()` inside `runCatching`. Returns `Result<String>` — the public URL on success, a wrapped exception on failure.

### presign()
Serialises a `PresignRequest` to JSON, POSTs to `$backendBaseUrl/api/v1/storage/presign`, sets `Authorization: Bearer <idToken>`. Throws via `check` if the response is not successful or the body is empty.

### uploadToR2()
PUTs `fileBytes` to the presigned URL with the given `contentType`. Uses `ByteArray.toRequestBody(MediaType)`. Throws via `check` if the response is not successful.

## Dependencies
The implementation uses OkHttp for HTTP and `kotlinx-serialization-json` for JSON. Both must be declared in `gradle/libs.versions.toml` and added to `app/build.gradle.kts`.

## Testing
Inject a `MockWebServer` (OkHttp test library) instance and pass its URL as `backendBaseUrl` and an `OkHttpClient` pointed at it. Enqueue mock responses for the presign call and the R2 PUT call, then assert on `upload()` returning the expected public URL.
