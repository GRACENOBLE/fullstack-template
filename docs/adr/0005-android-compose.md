# ADR 0005 — Android + Kotlin + Jetpack Compose + Material3

**Status:** Accepted  
**Date:** 2026-06-25

## Context

The template needs a mobile layer. Key requirements:
- Native Android UI (not a WebView wrapper)
- Type-safe, concise UI code that co-locates layout and logic
- Consistent with Google's current design system recommendations
- Single-activity architecture to avoid Fragment back-stack complexity

The template targets Android only; iOS would require a separate React Native or Flutter track.

## Decision

Use **Android** with **Kotlin 2.2**, **Jetpack Compose BOM 2026.02**, and **Material3**.

Key conventions enforced:
- **Single Activity** — `MainActivity` is the sole Android Activity. All navigation is handled by `Navigation Compose` and `AppNavGraph`.
- **No logic in `@Composable` functions** — screens are stateless; state is hoisted to ViewModels or the calling composable.
- **Material3 only** — `MaterialTheme.colorScheme.*` and `MaterialTheme.typography.*` for all colours and text styles; no hardcoded values.
- **Version catalog** — all dependency versions declared in `mobile/gradle/libs.versions.toml`.
- **`UiState<T>` sealed class** — generic loading/error/success state for ViewModels that fetch data.
- **Spotless + ktlint** — code formatting enforced in CI via `./gradlew spotlessCheck`.

## Consequences

### Positive
- Compose's declarative model eliminates XML layout files and `findViewById` boilerplate.
- Material3 theming tokens mean a single colour palette change propagates everywhere.
- The Single Activity / Navigation Compose pattern makes deep links and back-stack management straightforward.
- Kotlin coroutines + `StateFlow` integrate cleanly with Compose's `collectAsStateWithLifecycle()`.
- JVM unit tests (MockWebServer, JUnit 4) run without an emulator; instrumented tests (`createComposeRule()`) verify UI on device/emulator.

### Negative / trade-offs
- iOS is not covered. Adding iOS requires a separate codebase (Swift/SwiftUI) or migrating to a cross-platform framework.
- Compose's recomposition model requires care around `remember`, `derivedStateOf`, and `LaunchedEffect` — misuse causes performance issues or stale state.
- Instrumented tests require a running emulator or physical device and are excluded from the standard CI gate.
- `google-services.json` must be written to `mobile/app/` during CI from a secret; it is gitignored locally.
