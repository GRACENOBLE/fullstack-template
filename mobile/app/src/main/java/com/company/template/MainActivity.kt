package com.company.template

import android.Manifest
import android.annotation.SuppressLint
import android.os.Build
import android.os.Bundle
import android.util.Log
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.activity.result.contract.ActivityResultContracts
import androidx.activity.viewModels
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import com.company.template.auth.AuthViewModel
import com.company.template.auth.FirebaseAuthRepository
import com.company.template.navigation.AppNavGraph
import com.company.template.navigation.AppViewModel
import com.company.template.onboarding.DataStoreOnboardingRepository
import com.company.template.ui.theme.TemplateTheme
import com.google.firebase.messaging.FirebaseMessaging
import io.sentry.android.core.SentryAndroid

/**
 * Returns true when [dsn] is non-blank — i.e. Sentry should be initialised.
 * Extracted as a top-level function so it can be unit-tested on the JVM
 * without touching any Android framework APIs.
 */
fun shouldInitSentry(dsn: String): Boolean = dsn.isNotBlank()

class MainActivity : ComponentActivity() {

    // activity-compose 1.8.0 transitively pulls in Fragment 1.6+; lint can't detect this
    @SuppressLint("InvalidFragmentVersionForActivityResult")
    private val requestNotificationPermission =
        registerForActivityResult(ActivityResultContracts.RequestPermission()) { /* no-op */ }

    private val authRepository by lazy { FirebaseAuthRepository() }
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

@Preview(showBackground = true)
@Composable
fun AppPreview() {
    TemplateTheme {
        Greeting("Android")
    }
}

@Composable
fun Greeting(name: String, modifier: Modifier = Modifier) {
    androidx.compose.material3.Text(
        text = "Hello $name!",
        modifier = modifier
    )
}
