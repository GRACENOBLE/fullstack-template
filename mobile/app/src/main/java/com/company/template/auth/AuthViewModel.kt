package com.company.template.auth

import android.app.Activity
import android.content.Context
import android.content.Intent
import android.provider.Settings
import androidx.credentials.CredentialManager
import androidx.credentials.CustomCredential
import androidx.credentials.GetCredentialRequest
import androidx.credentials.exceptions.NoCredentialException
import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import androidx.lifecycle.viewModelScope
import androidx.lifecycle.viewmodel.initializer
import androidx.lifecycle.viewmodel.viewModelFactory
import com.company.template.BuildConfig
import com.google.android.libraries.identity.googleid.GetSignInWithGoogleOption
import com.google.android.libraries.identity.googleid.GoogleIdTokenCredential
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch

sealed class AuthUiState {
    data object Idle : AuthUiState()

    data object Loading : AuthUiState()

    data object Success : AuthUiState()

    data class Error(
        val message: String,
    ) : AuthUiState()
}

data class LoginFormState(
    val email: String = "",
    val password: String = "",
)

data class RegisterFormState(
    val name: String = "",
    val email: String = "",
    val password: String = "",
    val confirmPassword: String = "",
)

class AuthViewModel(
    private val repo: AuthRepository,
) : ViewModel() {
    private val _uiState = MutableStateFlow<AuthUiState>(AuthUiState.Idle)
    val uiState: StateFlow<AuthUiState> = _uiState.asStateFlow()

    val currentUser: StateFlow<User?> = repo.authStateFlow

    private val _loginForm = MutableStateFlow(LoginFormState())
    val loginForm: StateFlow<LoginFormState> = _loginForm.asStateFlow()

    private val _registerForm = MutableStateFlow(RegisterFormState())
    val registerForm: StateFlow<RegisterFormState> = _registerForm.asStateFlow()

    // Login form updates — clear error on any field change
    fun updateLoginEmail(email: String) {
        _loginForm.update { it.copy(email = email) }
        clearError()
    }

    fun updateLoginPassword(password: String) {
        _loginForm.update { it.copy(password = password) }
        clearError()
    }

    // Register form updates
    fun updateRegisterName(name: String) {
        _registerForm.update { it.copy(name = name) }
        clearError()
    }

    fun updateRegisterEmail(email: String) {
        _registerForm.update { it.copy(email = email) }
        clearError()
    }

    fun updateRegisterPassword(password: String) {
        _registerForm.update { it.copy(password = password) }
        clearError()
    }

    fun updateRegisterConfirmPassword(confirmPassword: String) {
        _registerForm.update { it.copy(confirmPassword = confirmPassword) }
        clearError()
    }

    fun signIn() {
        val (email, password) = _loginForm.value
        viewModelScope.launch {
            _uiState.value = AuthUiState.Loading
            repo
                .signInWithEmail(email, password)
                .onSuccess { _uiState.value = AuthUiState.Success }
                .onFailure { _uiState.value = AuthUiState.Error(it.message ?: "Sign in failed") }
        }
    }

    fun register() {
        val (name, email, password, confirm) = _registerForm.value
        if (password != confirm) {
            _uiState.value = AuthUiState.Error("Passwords do not match")
            return
        }
        viewModelScope.launch {
            _uiState.value = AuthUiState.Loading
            repo
                .registerWithEmail(name, email, password)
                .onSuccess { _uiState.value = AuthUiState.Success }
                .onFailure { _uiState.value = AuthUiState.Error(it.message ?: "Registration failed") }
        }
    }

    fun signInWithGoogle(activity: Activity) {
        viewModelScope.launch {
            _uiState.value = AuthUiState.Loading
            fetchGoogleIdToken(activity)
                .onSuccess { token ->
                    repo
                        .signInWithGoogle(token)
                        .onSuccess { _uiState.value = AuthUiState.Success }
                        .onFailure { _uiState.value = AuthUiState.Error(it.message ?: "Sign in failed") }
                }.onFailure { _uiState.value = AuthUiState.Error(it.message ?: "Google sign in failed") }
        }
    }

    fun signOut() {
        viewModelScope.launch {
            repo.signOut()
            _uiState.value = AuthUiState.Idle
        }
    }

    fun clearError() {
        if (_uiState.value is AuthUiState.Error) _uiState.value = AuthUiState.Idle
    }

    private suspend fun fetchGoogleIdToken(activity: Activity): Result<String> =
        runCatching {
            val webClientId = resolveWebClientId(activity)
            check(webClientId.isNotEmpty()) {
                "Google Sign-In not configured. Enable Google Sign-In in Firebase Console and " +
                    "re-download google-services.json, or add GOOGLE_WEB_CLIENT_ID to local.properties."
            }

            val credentialManager = CredentialManager.create(activity)
            val option = GetSignInWithGoogleOption.Builder(webClientId).build()
            val request = GetCredentialRequest.Builder().addCredentialOption(option).build()

            val result =
                try {
                    credentialManager.getCredential(activity, request)
                } catch (e: NoCredentialException) {
                    activity.startActivity(
                        Intent(Settings.ACTION_ADD_ACCOUNT).apply {
                            putExtra(Settings.EXTRA_ACCOUNT_TYPES, arrayOf("com.google"))
                        },
                    )
                    error("No Google account found on device. Please add an account and try again.")
                }

            val credential = result.credential
            check(
                credential is CustomCredential &&
                    credential.type == GoogleIdTokenCredential.TYPE_GOOGLE_ID_TOKEN_CREDENTIAL,
            ) { "Unexpected credential type: ${credential.type}" }

            GoogleIdTokenCredential.createFrom(credential.data).idToken
        }

    private fun resolveWebClientId(context: Context): String {
        val resId =
            context.resources.getIdentifier(
                "default_web_client_id",
                "string",
                context.packageName,
            )
        if (resId != 0) {
            val fromResource = context.getString(resId)
            if (fromResource.isNotEmpty()) return fromResource
        }
        return BuildConfig.GOOGLE_WEB_CLIENT_ID
    }

    companion object {
        fun factory(repo: AuthRepository): ViewModelProvider.Factory =
            viewModelFactory {
                initializer { AuthViewModel(repo) }
            }
    }
}
