---
topic: Activity and Compose architecture
last_verified: 2026-06-25
sources:
  - app/src/main/java/com/company/template/MainActivity.kt
  - app/build.gradle.kts
  - gradle/libs.versions.toml
  - app/src/main/java/com/company/template/data/network/ApiClient.kt
  - app/src/main/java/com/company/template/data/network/ApiResponse.kt
  - app/src/main/java/com/company/template/data/network/UserApi.kt
  - app/src/main/java/com/company/template/auth/AuthViewModel.kt
  - app/src/main/java/com/company/template/auth/LoginScreen.kt
  - app/src/main/java/com/company/template/auth/RegisterScreen.kt
  - app/src/main/java/com/company/template/home/HomeScreen.kt
  - app/src/main/java/com/company/template/home/HomeViewModel.kt
  - app/src/main/java/com/company/template/navigation/AppNavGraph.kt
  - app/src/main/java/com/company/template/ui/state/UiState.kt
  - app/src/main/java/com/company/template/ui/components/UiStateContent.kt
---

# Activity and Compose architecture

## Entry point

`MainActivity` is the single Activity for the app. It extends `ComponentActivity` and sets up the Compose UI tree in `onCreate`:

```kotlin
class MainActivity : ComponentActivity() {

    @SuppressLint("InvalidFragmentVersionForActivityResult")
    private val requestNotificationPermission =
        registerForActivityResult(ActivityResultContracts.RequestPermission()) { /* no-op */ }

    private val authRepository by lazy { FirebaseAuthRepository(applicationContext) }
    private val onboardingRepository by lazy { DataStoreOnboardingRepository(applicationContext) }

    private val authViewModel: AuthViewModel by viewModels {
        AuthViewModel.factory(authRepository)
    }
    private val appViewModel: AppViewModel by viewModels {
        AppViewModel.factory(authRepository, onboardingRepository)
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            requestNotificationPermission.launch(Manifest.permission.POST_NOTIFICATIONS)
        }
        if (shouldInitSentry(BuildConfig.SENTRY_DSN)) {
            SentryAndroid.init(this) { options ->
                options.dsn = BuildConfig.SENTRY_DSN
                options.tracesSampleRate = 1.0
            }
        }
        FirebaseMessaging.getInstance().token.addOnSuccessListener { token ->
            Log.d("FCM_TOKEN", token)
        }
        enableEdgeToEdge()
        setContent {
            TemplateTheme {
                Scaffold(modifier = Modifier.fillMaxSize()) { innerPadding ->
                    AppNavGraph(
                        appViewModel = appViewModel,
                        authViewModel = authViewModel,
                        modifier = Modifier.padding(innerPadding)
                    )
                }
            }
        }
    }
}
```

Key calls:
- `requestNotificationPermission` — registered at construction time; launched on API 33+ (TIRAMISU) to request `POST_NOTIFICATIONS` at runtime.
- `shouldInitSentry(dsn)` — top-level function (unit-testable on JVM) that returns `true` when `BuildConfig.SENTRY_DSN` is non-blank; Sentry is only initialised when a DSN is present.
- `FirebaseMessaging.getInstance().token` — fetches the FCM registration token and logs it at startup.
- `enableEdgeToEdge()` — called before `setContent`; allows content to draw behind system bars.
- `setContent { }` — replaces the XML layout system; the lambda is the Compose UI root.
- `TemplateTheme { }` — applied once here; all screens inherit the theme automatically.
- `AppNavGraph` — the top-level navigation graph; receives `appViewModel` and `authViewModel` as parameters.

## ViewModels wired in MainActivity

`AuthViewModel` and `AppViewModel` are both scoped to the Activity via `by viewModels { }` with custom factories:

- `AuthViewModel.factory(authRepository)` — owns Firebase auth state.
- `AppViewModel.factory(authRepository, onboardingRepository)` — drives the initial navigation decision (onboarding vs. home).

Both repositories are instantiated lazily with `applicationContext` to avoid Activity leaks.

## Feature-scoped ViewModel state: AuthUiState

`AuthViewModel` exposes its operation state as a `StateFlow<AuthUiState>` where `AuthUiState` is a sealed class defined in `auth/AuthViewModel.kt`:

```kotlin
sealed class AuthUiState {
    data object Idle : AuthUiState()
    data object Loading : AuthUiState()
    data object Success : AuthUiState()
    data class Error(val message: String) : AuthUiState()
}
```

- `Idle` — no operation in progress; initial state and the state restored after sign-out or error dismissal.
- `Loading` — an async auth operation (sign-in, register, Google sign-in) is in flight.
- `Success` — the operation completed successfully; `AppNavGraph` observes this to navigate to the home screen.
- `Error(message)` — the operation failed; `message` is displayed inline in the auth screen.

`AuthViewModel` also exposes two separate form state flows — `loginForm: StateFlow<LoginFormState>` and `registerForm: StateFlow<RegisterFormState>` — so each screen receives only its own form fields. Any field update clears an active `Error` state via `clearError()`.

This is a feature-scoped sealed class suited to auth operations where `Success` carries no payload. For ViewModels that load data and expose it to the UI, use the generic `UiState<T>` sealed class instead — see `mobile/docs/ui-states.md`.

## Generic UI state: UiState\<T\>

`ui/state/UiState.kt` defines a generic sealed class for ViewModels that load typed data:

```kotlin
sealed class UiState<out T> {
    data object Idle : UiState<Nothing>()
    data object Loading : UiState<Nothing>()
    data class Success<T>(val data: T) : UiState<T>()
    data class Error(val message: String) : UiState<Nothing>()
}
```

`HomeViewModel` uses this pattern: it exposes `profileState: StateFlow<UiState<UserProfile>>` and provides a `refresh()` method that re-runs the fetch. `HomeScreen` collects the flow with `collectAsStateWithLifecycle()` and delegates rendering to `UiStateContent`.

See `mobile/docs/ui-states.md` for the full pattern, wiring instructions, and testing approach.

## Single-Activity pattern

There is no `Fragment` stack. All navigation between screens happens inside the Compose composition via Navigation Compose. Do not create additional Activities or Fragments.

## Lifecycle

`ComponentActivity` integrates with Jetpack Lifecycle. When adding state or coroutines:
- Use `viewModel()` (from `lifecycle-viewmodel-compose`) to scope ViewModels to the Activity or a nav destination.
- Collect `StateFlow` / `Flow` from ViewModels using `collectAsStateWithLifecycle()` (from `lifecycle-runtime-compose`) — not `collectAsState()`, which does not respect lifecycle.

## Build configuration

- `compileSdk` 36 (minor API level 1), `minSdk` 24, `targetSdk` 36
- Source and target compatibility: Java 11
- `buildFeatures { compose = true; buildConfig = true }` — enables the Compose compiler and the `BuildConfig` class generation
- Kotlin plugin: `org.jetbrains.kotlin.plugin.compose` and `org.jetbrains.kotlin.plugin.serialization`
- Google Services plugin: `com.google.gms.google-services`

### buildConfigField entries

Both fields are read from `local.properties` at build time. They default to an empty string if the key is absent:

| Field | Type | Source key |
|---|---|---|
| `SENTRY_DSN` | `String` | `local.properties: SENTRY_DSN` |
| `GOOGLE_WEB_CLIENT_ID` | `String` | `local.properties: GOOGLE_WEB_CLIENT_ID` |

Release build validation: assembling or bundling a release variant fails fast if `GOOGLE_WEB_CLIENT_ID` is empty in `local.properties`.

| Field | Type | Source key | Default |
|---|---|---|---|
| `BACKEND_URL` | `String` | `local.properties: BACKEND_URL` | `http://10.0.2.2:8080` |

`http://10.0.2.2:8080` is the Android emulator alias for host `localhost`. Override in `local.properties` for physical devices or staging environments.

## Dependency versions

All versions are declared in `gradle/libs.versions.toml`. Current versions:

| Library / group | Version |
|---|---|
| Kotlin | 2.2.10 |
| AGP | 9.2.1 |
| Google Services plugin | 4.4.2 |
| Compose BOM | 2026.02.01 |
| `androidx.core:core-ktx` | 1.10.1 |
| `androidx.lifecycle:lifecycle-runtime-ktx` | 2.6.1 |
| `androidx.lifecycle:lifecycle-viewmodel-compose` | 2.9.0 |
| `androidx.lifecycle:lifecycle-runtime-compose` | 2.9.0 |
| `androidx.activity:activity-compose` | 1.8.0 |
| `androidx.navigation:navigation-compose` | 2.9.0 |
| `androidx.credentials:credentials` | 1.5.0 |
| `androidx.credentials:credentials-play-services-auth` | 1.5.0 |
| `androidx.datastore:datastore-preferences` | 1.1.7 |
| `com.google.android.libraries.identity.googleid:googleid` | 1.1.1 |
| `com.google.firebase:firebase-bom` | 33.7.0 |
| `io.sentry:sentry-android` | 8.14.0 |
| `com.squareup.okhttp3:okhttp` | 4.12.0 |
| `org.jetbrains.kotlinx:kotlinx-coroutines-android` | 1.10.2 |
| `org.jetbrains.kotlinx:kotlinx-serialization-json` | 1.8.1 |
| `io.coil-kt:coil-compose` | 2.7.0 |

Compose library versions (ui, material3, etc.) are managed by the BOM — do not pin them individually. Firebase library versions (firebase-messaging-ktx, firebase-auth-ktx, etc.) are managed by the Firebase BOM — do not pin them individually.

## Data layer

The `data/network/` package contains the HTTP client and all API call functions. Key types:

- `ApiClient` — singleton `OkHttpClient` with `AuthInterceptor` that attaches the Firebase ID token as `Authorization: Bearer <token>` on every request.
- `ApiResponse<T>` / `ApiErrorResponse` / `ApiErrorDetail` — envelope types that match the Go backend's `{"data": ...}` / `{"error": {"code": "...", "message": "..."}}` response shape.
- `UserApi` — suspend functions for the `/api/v1/me` endpoint; accepts injectable `baseUrl` and `client` parameters for unit testing with `MockWebServer`.

See `mobile/docs/http-client.md` for the full pattern, instructions on adding new API calls, and the testing approach.
