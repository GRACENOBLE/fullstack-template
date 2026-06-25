---
topic: index
last_verified: 2026-06-25
---

# Mobile docs index

Topic-based documentation for the Android app. Each file is kept in sync with the source code it describes.

| Topic | File | Sources |
|---|---|---|
| Jetpack Compose UI conventions | `compose-conventions.md` | `app/src/main/java/com/company/template/MainActivity.kt`, `ui/theme/` |
| Activity and Compose architecture | `architecture.md` | `app/src/main/java/com/company/template/MainActivity.kt`, `app/build.gradle.kts` |
| Testing patterns | `testing.md` | `app/src/test/java/com/company/template/GreetingFormatTest.kt`, `app/src/androidTest/java/com/company/template/GreetingTest.kt`, `app/build.gradle.kts` |
| Observability (Sentry error tracking) | `observability.md` | `gradle/libs.versions.toml`, `app/build.gradle.kts`, `app/src/main/java/com/company/template/MainActivity.kt` |
| Firebase Cloud Messaging — service, token registration, background notifications | `fcm.md` | `app/src/main/java/com/company/template/fcm/MyFirebaseMessagingService.kt`, `app/src/main/java/com/company/template/fcm/FcmRegistrationPayload.kt`, `app/src/main/AndroidManifest.xml`, `gradle/libs.versions.toml` |
| Object storage (Cloudflare R2) — UploadRepository interface, R2UploadRepository, presign + PUT flow | `storage.md` | `app/src/main/java/com/company/template/storage/UploadRepository.kt` |
| HTTP client and API layer — ApiClient, envelope types, UserApi, MockWebServer testing, BACKEND_URL | `http-client.md` | `app/src/main/java/com/company/template/data/network/ApiClient.kt`, `ApiResponse.kt`, `UserApi.kt`, `app/src/test/java/com/company/template/data/network/UserApiTest.kt` |
| Generic UI state — UiState\<T\> sealed class, UiStateContent composable, ViewModel wiring, testing | `ui-states.md` | `app/src/main/java/com/company/template/ui/state/UiState.kt`, `ui/components/UiStateContent.kt`, `home/HomeViewModel.kt`, `home/HomeScreen.kt` |
