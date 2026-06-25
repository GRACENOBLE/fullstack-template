---
topic: HTTP client and API layer
last_verified: 2026-06-25
sources:
  - app/src/main/java/com/company/template/data/network/ApiClient.kt
  - app/src/main/java/com/company/template/data/network/ApiResponse.kt
  - app/src/main/java/com/company/template/data/network/UserApi.kt
  - app/src/main/java/com/company/template/home/HomeScreen.kt
  - app/src/test/java/com/company/template/data/network/UserApiTest.kt
  - app/build.gradle.kts
---

# HTTP client and API layer

## Package structure

```
data/network/
  ApiClient.kt      — singleton OkHttpClient with Firebase auth interceptor
  ApiResponse.kt    — envelope types matching the Go backend's response shape
  UserApi.kt        — UserProfile model + getMe() suspend function
```

## ApiClient

`ApiClient` is a Kotlin `object` (singleton). It exposes a single `httpClient: OkHttpClient` that has `AuthInterceptor` attached.

```kotlin
object ApiClient {
    val httpClient: OkHttpClient = OkHttpClient.Builder()
        .addInterceptor(AuthInterceptor())
        .connectTimeout(30, TimeUnit.SECONDS)
        .readTimeout(30, TimeUnit.SECONDS)
        .build()
}
```

### AuthInterceptor

`AuthInterceptor` is a private inner class that adds `Authorization: Bearer <token>` to every outgoing request.

```kotlin
private class AuthInterceptor : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val token = runCatching {
            FirebaseAuth.getInstance().currentUser
                ?.getIdToken(false)
                ?.result
                ?.token
        }.getOrNull()

        val request = if (token != null) {
            chain.request().newBuilder()
                .header("Authorization", "Bearer $token")
                .build()
        } else {
            chain.request()
        }
        return chain.proceed(request)
    }
}
```

`getIdToken(false)` returns the cached token if it is still valid; it does not force a refresh. The `.result` property is the synchronous accessor on the Firebase `Task<T>`. This is safe here because `intercept()` always executes on an OkHttp dispatcher thread, never on the main thread.

If the current user is null or token retrieval fails, the request is forwarded without an `Authorization` header rather than throwing.

## Envelope types

The Go backend wraps all responses in a consistent envelope. The Kotlin types mirror that shape exactly:

| Go shape | Kotlin type |
|---|---|
| `{"data": <T>}` | `ApiResponse<T>(val data: T)` |
| `{"error": {"code": "...", "message": "..."}}` | `ApiErrorResponse(val error: ApiErrorDetail)` |
| `{"code": "...", "message": "..."}` | `ApiErrorDetail(val code: String, val message: String)` |

All three types are annotated with `@Serializable` (kotlinx.serialization).

## UserApi

`UserApi` is a Kotlin `object` with a single suspend function `getMe()` that calls `GET /api/v1/me` and returns `Result<UserProfile>`.

```kotlin
object UserApi {
    private val json = Json { ignoreUnknownKeys = true }

    suspend fun getMe(
        baseUrl: String = BuildConfig.BACKEND_URL,
        client: OkHttpClient = ApiClient.httpClient,
    ): Result<UserProfile> = runCatching {
        val request = Request.Builder()
            .url("$baseUrl/api/v1/me")
            .get()
            .build()

        client.newCall(request).execute().use { response ->
            val body = response.body?.string() ?: error("empty body")
            if (!response.isSuccessful) {
                val err = json.decodeFromString<ApiErrorResponse>(body)
                error(err.error.message)
            }
            json.decodeFromString<ApiResponse<UserProfile>>(body).data
        }
    }
}
```

`UserProfile` is defined in the same file:

```kotlin
@Serializable
data class UserProfile(
    val uid: String,
    val email: String? = null,
    val displayName: String? = null,
)
```

`ignoreUnknownKeys = true` on the `Json` instance ensures the client does not break when the backend adds new fields.

## BACKEND_URL build config field

`BuildConfig.BACKEND_URL` is injected at build time from `local.properties`:

```kotlin
// app/build.gradle.kts
buildConfigField(
    "String",
    "BACKEND_URL",
    "\"${localProps.getProperty("BACKEND_URL", "http://10.0.2.2:8080")}\"",
)
```

`http://10.0.2.2:8080` is the Android emulator's alias for the host machine's `localhost`. This default works for local development without any `local.properties` entry.

To override for a physical device or a staging server, add to `mobile/local.properties`:

```properties
BACKEND_URL=http://192.168.1.x:8080
```

`local.properties` is gitignored; never commit it.

## Calling an API from a Composable

Use `LaunchedEffect` to launch the coroutine and `withContext(Dispatchers.IO)` to move the blocking OkHttp call off the main thread. Store result in `remember` state:

```kotlin
var profile by remember { mutableStateOf<UserProfile?>(null) }
var profileError by remember { mutableStateOf<String?>(null) }

LaunchedEffect(Unit) {
    val result = withContext(Dispatchers.IO) { UserApi.getMe() }
    result
        .onSuccess { profile = it }
        .onFailure { profileError = it.message }
}
```

Source: `home/HomeScreen.kt`.

## Adding a new API call

Follow the `UserApi` pattern:

1. Add a new `@Serializable` response model in `data/network/`.
2. Add a suspend function to an `object` (or a new `object` for a new resource) with `baseUrl` and `client` as defaulted parameters.
3. Use `runCatching { }` so all exceptions are captured in `Result`.
4. Decode success bodies as `ApiResponse<YourModel>`, error bodies as `ApiErrorResponse`.
5. Call from a Composable via `LaunchedEffect` + `withContext(Dispatchers.IO)`.

## Testing

API functions are tested with `MockWebServer` (from `com.squareup.okhttp3:mockwebserver`) in JVM unit tests under `src/test/`. No Android framework is needed.

```kotlin
class UserApiTest {

    private lateinit var server: MockWebServer
    private val testClient = OkHttpClient()   // no AuthInterceptor

    @Before fun setUp() { server = MockWebServer(); server.start() }
    @After  fun tearDown() { server.shutdown() }

    @Test
    fun `getMe returns UserProfile on successful response`() = runTest {
        server.enqueue(
            MockResponse()
                .setResponseCode(200)
                .setBody("""{"data":{"uid":"u1","email":"a@b.com"}}""")
        )

        val result = UserApi.getMe(
            baseUrl = server.url("/").toString().trimEnd('/'),
            client = testClient,
        )

        assertTrue(result.isSuccess)
        assertEquals("u1", result.getOrThrow().uid)
    }

    @Test
    fun `getMe returns failure with backend message on error response`() = runTest {
        server.enqueue(
            MockResponse()
                .setResponseCode(401)
                .setBody("""{"error":{"code":"UNAUTHENTICATED","message":"no token"}}""")
        )

        val result = UserApi.getMe(
            baseUrl = server.url("/").toString().trimEnd('/'),
            client = testClient,
        )

        assertTrue(result.isFailure)
        assertEquals("no token", result.exceptionOrNull()?.message)
    }
}
```

Key points:
- Pass the `MockWebServer` URL as `baseUrl` and a plain `OkHttpClient()` as `client` — this bypasses `AuthInterceptor` and avoids any Firebase dependency in tests.
- Use `kotlinx-coroutines-test` (`runTest`) for suspend functions.
- Tests live in `src/test/` (JVM), not `src/androidTest/`.
