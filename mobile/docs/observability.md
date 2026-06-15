---
topic: observability
last_verified: 2026-06-15
sources:
  - gradle/libs.versions.toml
  - app/build.gradle.kts
  - app/src/main/java/com/company/template/MainActivity.kt
---

# Observability

## Sentry SDK

| Dependency | Version |
|---|---|
| `io.sentry:sentry-android` | 8.14.0 |

Declared in `gradle/libs.versions.toml`:

```toml
[versions]
sentry = "8.14.0"

[libraries]
sentry-android = { group = "io.sentry", name = "sentry-android", version.ref = "sentry" }
```

Added to `app/build.gradle.kts`:

```kotlin
implementation(libs.sentry.android)
```

## BuildConfig.SENTRY_DSN

`app/build.gradle.kts` sets `buildConfig = true` and defines the field in `defaultConfig`:

```kotlin
buildFeatures {
    buildConfig = true
}

defaultConfig {
    buildConfigField("String", "SENTRY_DSN", "\"\"")
}
```

The default value is an empty string, meaning Sentry is off for all build variants unless overridden.

**To supply a real DSN — per build type:**

```kotlin
buildTypes {
    release {
        buildConfigField("String", "SENTRY_DSN", "\"https://<key>@o<org>.ingest.sentry.io/<project>\"")
    }
}
```

**To supply a real DSN — via command line:**

```bash
./gradlew assembleRelease -PSENTRY_DSN="https://<key>@o<org>.ingest.sentry.io/<project>"
```

Then reference the project property in `build.gradle.kts`:

```kotlin
buildConfigField("String", "SENTRY_DSN", "\"${project.findProperty("SENTRY_DSN") ?: ""}\"")
```

## shouldInitSentry helper

`MainActivity.kt` defines a top-level function:

```kotlin
fun shouldInitSentry(dsn: String): Boolean = dsn.isNotBlank()
```

It is a pure function with no Android framework dependencies, making it directly unit-testable on the JVM without Robolectric or an emulator.

## Initialization in MainActivity.onCreate

`SentryAndroid.init` is called before `setContent`, guarded by `shouldInitSentry`:

```kotlin
override fun onCreate(savedInstanceState: Bundle?) {
    super.onCreate(savedInstanceState)
    if (shouldInitSentry(BuildConfig.SENTRY_DSN)) {
        SentryAndroid.init(this) { options ->
            options.dsn = BuildConfig.SENTRY_DSN
            options.tracesSampleRate = 1.0
        }
    }
    enableEdgeToEdge()
    setContent { ... }
}
```

When `BuildConfig.SENTRY_DSN` is blank (the default), the `if` block is skipped entirely and no Sentry code runs.
