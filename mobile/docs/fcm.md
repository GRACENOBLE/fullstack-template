---
topic: fcm
last_verified: 2026-06-16
sources:
  - app/src/main/java/com/company/template/fcm/FcmRegistrationPayload.kt
  - app/src/main/java/com/company/template/fcm/MyFirebaseMessagingService.kt
  - app/src/main/AndroidManifest.xml
  - app/src/test/java/com/company/template/fcm/FcmRegistrationPayloadTest.kt
  - gradle/libs.versions.toml
  - app/build.gradle.kts
---

# Firebase Cloud Messaging (FCM) — Android

Handles foreground and background push notifications via `FirebaseMessagingService`.

## Dependencies

`gradle/libs.versions.toml`:
```toml
[versions]
firebaseBom = "33.7.0"

[libraries]
firebase-bom = { group = "com.google.firebase", name = "firebase-bom", version.ref = "firebaseBom" }
firebase-messaging-ktx = { group = "com.google.firebase", name = "firebase-messaging-ktx" }
```

`app/build.gradle.kts`:
```kotlin
implementation(platform(libs.firebase.bom))
implementation(libs.firebase.messaging.ktx)
```

> **Note:** The Google Services Gradle plugin (`com.google.gms.google-services`) and a `google-services.json` file are required for Firebase to initialise at runtime. Add the plugin to `build.gradle.kts` and place `google-services.json` in `app/` when setting up a real Firebase project. Without them the app compiles but Firebase will not initialise.

## AndroidManifest.xml

```xml
<uses-permission android:name="android.permission.POST_NOTIFICATIONS" />
<uses-permission android:name="android.permission.INTERNET" />

<service
    android:name=".fcm.MyFirebaseMessagingService"
    android:exported="false">
    <intent-filter>
        <action android:name="com.google.firebase.MESSAGING_EVENT" />
    </intent-filter>
</service>
```

`POST_NOTIFICATIONS` is required on Android 13+ (API 33+) for runtime notification permission. `INTERNET` is required to POST the token to the backend.

## FcmRegistrationPayload

`com.company.template.fcm.FcmRegistrationPayload` is a `@Serializable` data class representing the JSON body sent to `POST /api/v1/fcm/register`:

```kotlin
@Serializable
data class FcmRegistrationPayload(
    val token: String,
    val platform: String = "android",
)
```

Serialised with `kotlinx.serialization.json.Json`. Platform defaults to `"android"`.

## MyFirebaseMessagingService

`com.company.template.fcm.MyFirebaseMessagingService` extends `FirebaseMessagingService`.

### onNewToken

Called by the FCM SDK when a new or refreshed token is available. Launches a coroutine on `Dispatchers.IO` to POST the token to the backend:

1. Reads `backend_base_url` from string resources (define in `res/values/strings.xml` for each build variant).
2. Serialises a `FcmRegistrationPayload` and sends it with OkHttp.
3. Failures are swallowed with `runCatching` — token registration is best-effort.

```kotlin
override fun onNewToken(token: String) {
    serviceScope.launch { registerTokenWithBackend(token) }
}
```

### onMessageReceived

Called for foreground messages (app in foreground). Shows a system notification via `NotificationCompat.Builder`.

`NotificationChannel` creation is guarded by `Build.VERSION.SDK_INT >= Build.VERSION_CODES.O` because `NotificationChannel` was added in API 26 and `minSdk = 24`.

Background messages (app not running or in background) are handled automatically by the FCM SDK without calling `onMessageReceived`.

## Configuring the backend URL

Add `backend_base_url` to `app/src/main/res/values/strings.xml`:
```xml
<string name="backend_base_url">https://api.example.com</string>
```

Override per build type in `app/src/debug/res/values/strings.xml`:
```xml
<string name="backend_base_url">http://10.0.2.2:8080</string>
```

(`10.0.2.2` is the Android emulator's alias for the host machine's `localhost`.)

## Testing

Unit tests in `app/src/test/java/com/company/template/fcm/FcmRegistrationPayloadTest.kt` (JUnit 4, JVM):

- Verifies JSON serialisation contains the correct field names and values
- Verifies the default platform is `"android"`
- Verifies the platform can be overridden
- Verifies round-trip serialisation (encode → decode → equal)

No Android framework is needed for these tests. Run with `./gradlew test`.
