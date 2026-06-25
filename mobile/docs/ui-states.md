---
topic: Generic UI state (UiState<T> and UiStateContent)
last_verified: 2026-06-25
sources:
  - app/src/main/java/com/company/template/ui/state/UiState.kt
  - app/src/main/java/com/company/template/ui/components/UiStateContent.kt
  - app/src/main/java/com/company/template/home/HomeViewModel.kt
  - app/src/main/java/com/company/template/home/HomeScreen.kt
  - app/src/test/java/com/company/template/home/HomeViewModelTest.kt
  - app/src/androidTest/java/com/company/template/ui/components/UiStateContentTest.kt
---

# Generic UI state: UiState\<T\> and UiStateContent

## UiState\<T\> sealed class

`ui/state/UiState.kt` defines a covariant sealed class for ViewModels that load typed data:

```kotlin
sealed class UiState<out T> {
    data object Idle    : UiState<Nothing>()
    data object Loading : UiState<Nothing>()
    data class  Success<T>(val data: T) : UiState<T>()
    data class  Error(val message: String) : UiState<Nothing>()
}
```

| Variant | Meaning | Carries |
|---|---|---|
| `Idle` | No operation started yet; rendered as nothing | — |
| `Loading` | Fetch in progress | — |
| `Success<T>` | Fetch completed; data is ready | `data: T` |
| `Error` | Fetch failed | `message: String` |

Use `UiState<T>` for any ViewModel that loads a specific type from the network or a repository. For auth operations where `Success` has no payload, use the feature-scoped `AuthUiState` in `auth/AuthViewModel.kt` instead.

## UiStateContent composable

`ui/components/UiStateContent.kt` is the canonical rendering component for `UiState<T>`:

```kotlin
@Composable
fun <T> UiStateContent(
    state: UiState<T>,
    modifier: Modifier = Modifier,
    onRetry: (() -> Unit)? = null,
    content: @Composable (T) -> Unit,
)
```

| Parameter | Type | Default | Purpose |
|---|---|---|---|
| `state` | `UiState<T>` | — | State to render |
| `modifier` | `Modifier` | `Modifier` | Applied to the container `Box` for Loading/Error states |
| `onRetry` | `(() -> Unit)?` | `null` | When non-null, shows a "Retry" `TextButton` in the Error state |
| `content` | `@Composable (T) -> Unit` | — | Trailing lambda called with `state.data` when state is `Success` |

Rendering per variant:

- `Idle` — renders nothing (`Unit`).
- `Loading` — centered `CircularProgressIndicator` inside a full-width `Box` with `32.dp` padding. Test tag: `"loading_indicator"`.
- `Error` — centered error `Text` (color `MaterialTheme.colorScheme.error`, style `bodyMedium`) inside a full-width `Box` with `16.dp` padding. If `onRetry != null`, a `TextButton("Retry")` appears below the message.
- `Success` — calls `content(state.data)`; the component itself renders nothing else.

## Wiring in a ViewModel

`HomeViewModel` demonstrates the full pattern:

```kotlin
class HomeViewModel(
    private val baseUrl: String,
    private val httpClient: OkHttpClient,
    private val ioDispatcher: CoroutineDispatcher = Dispatchers.IO,
) : ViewModel() {

    private val _profileState = MutableStateFlow<UiState<UserProfile>>(UiState.Idle)
    val profileState: StateFlow<UiState<UserProfile>> = _profileState.asStateFlow()

    init { fetchProfile() }

    fun refresh() { fetchProfile() }

    private fun fetchProfile() {
        viewModelScope.launch(ioDispatcher) {
            _profileState.value = UiState.Loading
            UserApi.getMe(baseUrl = baseUrl, client = httpClient)
                .onSuccess { _profileState.value = UiState.Success(it) }
                .onFailure { _profileState.value = UiState.Error(it.message ?: "Failed to load profile") }
        }
    }

    companion object {
        fun factory(): ViewModelProvider.Factory = viewModelFactory {
            initializer {
                HomeViewModel(
                    baseUrl = BuildConfig.BACKEND_URL,
                    httpClient = ApiClient.httpClient,
                )
            }
        }
    }
}
```

Key points:
- Initial state is `Idle`; `fetchProfile()` is called in `init` so the fetch starts immediately.
- `ioDispatcher` is injected (defaults to `Dispatchers.IO`) so unit tests can substitute `UnconfinedTestDispatcher`.
- `baseUrl` and `httpClient` are constructor parameters so `MockWebServer` can be injected in tests.
- `refresh()` is a public method that re-runs the same fetch (sets state back to `Loading` first).

The factory uses the `viewModelFactory { initializer { … } }` DSL and reads `BuildConfig.BACKEND_URL` and `ApiClient.httpClient` for production use.

## Collecting in a screen

```kotlin
@Composable
fun HomeScreen(
    displayName: String,
    onSignOut: () -> Unit,
    viewModel: HomeViewModel = viewModel(factory = HomeViewModel.factory()),
    modifier: Modifier = Modifier,
) {
    val profileState by viewModel.profileState.collectAsStateWithLifecycle()

    UiStateContent(
        state = profileState,
        onRetry = viewModel::refresh,
    ) { profile ->
        ProfileContent(profile = profile)
    }
}
```

- Collect with `collectAsStateWithLifecycle()` (from `lifecycle-runtime-compose`), not `collectAsState()`.
- Pass `viewModel::refresh` as `onRetry` so the Retry button re-triggers the fetch.
- The screen contains no data-fetch logic — all transitions happen in the ViewModel.

## Testing

### ViewModel unit tests (JVM, `src/test/`)

`HomeViewModelTest` uses `MockWebServer` + `UnconfinedTestDispatcher`:

```kotlin
@OptIn(ExperimentalCoroutinesApi::class)
class HomeViewModelTest {
    private lateinit var server: MockWebServer
    private val testDispatcher = UnconfinedTestDispatcher()

    @Before fun setUp() {
        server = MockWebServer()
        server.start()
        Dispatchers.setMain(testDispatcher)
    }

    @After fun tearDown() {
        server.shutdown()
        Dispatchers.resetMain()
    }

    private fun buildViewModel() = HomeViewModel(
        baseUrl = server.url("/").toString().trimEnd('/'),
        httpClient = OkHttpClient(),
        ioDispatcher = testDispatcher,
    )
}
```

- Enqueue a `MockResponse` before constructing the ViewModel; `init` fires immediately with `UnconfinedTestDispatcher`.
- Assert on `viewModel.profileState.value` directly — no `runTest` `advanceUntilIdle()` required when using `UnconfinedTestDispatcher` with `runTest`.
- Test `refresh()` by enqueuing two responses: the first is consumed by `init`, the second by `refresh()`.

### UiStateContent instrumented tests (`src/androidTest/`)

`UiStateContentTest` uses `createComposeRule()` and sets each `UiState` variant in `setContent`:

```kotlin
@RunWith(AndroidJUnit4::class)
class UiStateContentTest {
    @get:Rule val composeTestRule = createComposeRule()

    @Test fun loadingState_showsCircularProgressIndicator() {
        composeTestRule.setContent {
            TemplateTheme {
                UiStateContent(state = UiState.Loading, content = { _: Unit -> })
            }
        }
        composeTestRule.onNodeWithTag("loading_indicator").assertIsDisplayed()
    }

    @Test fun errorState_withRetry_showsRetryButton() {
        var retryClicked = false
        composeTestRule.setContent {
            TemplateTheme {
                UiStateContent(
                    state = UiState.Error("Network error"),
                    onRetry = { retryClicked = true },
                    content = { _: Unit -> },
                )
            }
        }
        composeTestRule.onNodeWithText("Retry").performClick()
        assertTrue(retryClicked)
    }
}
```

- Wrap content in `TemplateTheme` so color tokens resolve.
- Use `onNodeWithTag("loading_indicator")` for the spinner (the tag is set in `UiStateContent`).
- Use `onNodeWithText(…)` for error messages and the Retry button label.
