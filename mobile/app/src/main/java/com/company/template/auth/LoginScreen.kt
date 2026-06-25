package com.company.template.auth

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.Button
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.Icon
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.SolidColor
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.graphics.vector.path
import androidx.compose.ui.platform.testTag
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.company.template.ui.theme.TemplateTheme

internal object LoginTestTags {
    const val SIGN_IN_BUTTON = "login_sign_in_button"
}

private val GoogleIcon: ImageVector
    get() =
        ImageVector
            .Builder(
                name = "GoogleIcon",
                defaultWidth = 24.dp,
                defaultHeight = 24.dp,
                viewportWidth = 24f,
                viewportHeight = 24f,
            ).apply {
                path(fill = SolidColor(Color(0xFF4285F4))) {
                    moveTo(22.56f, 12.25f)
                    curveToRelative(0f, -0.78f, -0.07f, -1.53f, -0.2f, -2.25f)
                    horizontalLineTo(12f)
                    verticalLineToRelative(4.26f)
                    horizontalLineToRelative(5.92f)
                    curveToRelative(-0.26f, 1.37f, -1.04f, 2.53f, -2.21f, 3.31f)
                    verticalLineToRelative(2.77f)
                    horizontalLineToRelative(3.57f)
                    curveToRelative(2.08f, -1.92f, 3.28f, -4.74f, 3.28f, -8.09f)
                    close()
                }
                path(fill = SolidColor(Color(0xFF34A853))) {
                    moveTo(12f, 23f)
                    curveToRelative(2.97f, 0f, 5.46f, -0.98f, 7.28f, -2.66f)
                    lineToRelative(-3.57f, -2.77f)
                    curveToRelative(-0.98f, 0.66f, -2.23f, 1.06f, -3.71f, 1.06f)
                    curveToRelative(-2.86f, 0f, -5.29f, -1.93f, -6.16f, -4.53f)
                    horizontalLineTo(2.18f)
                    verticalLineToRelative(2.84f)
                    curveTo(3.99f, 20.53f, 7.7f, 23f, 12f, 23f)
                    close()
                }
                path(fill = SolidColor(Color(0xFFFBBC05))) {
                    moveTo(5.84f, 14.09f)
                    curveToRelative(-0.22f, -0.66f, -0.35f, -1.36f, -0.35f, -2.09f)
                    reflectiveCurveToRelative(0.13f, -1.43f, 0.35f, -2.09f)
                    verticalLineTo(7.07f)
                    horizontalLineTo(2.18f)
                    curveTo(1.43f, 8.55f, 1f, 10.22f, 1f, 12f)
                    reflectiveCurveToRelative(0.43f, 3.45f, 1.18f, 4.93f)
                    lineToRelative(2.85f, -2.22f)
                    lineToRelative(0.81f, -0.62f)
                    close()
                }
                path(fill = SolidColor(Color(0xFFEA4335))) {
                    moveTo(12f, 5.38f)
                    curveToRelative(1.62f, 0f, 3.06f, 0.56f, 4.21f, 1.64f)
                    lineToRelative(3.15f, -3.15f)
                    curveTo(17.45f, 2.09f, 14.97f, 1f, 12f, 1f)
                    curveTo(7.7f, 1f, 3.99f, 3.47f, 2.18f, 7.07f)
                    lineToRelative(3.66f, 2.84f)
                    curveToRelative(0.87f, -2.6f, 3.3f, -4.53f, 6.16f, -4.53f)
                    close()
                }
            }.build()

@Composable
fun LoginScreen(
    loginForm: LoginFormState,
    uiState: AuthUiState,
    onEmailChange: (String) -> Unit,
    onPasswordChange: (String) -> Unit,
    onSignIn: () -> Unit,
    onSignInWithGoogle: () -> Unit,
    onNavigateToRegister: () -> Unit,
    modifier: Modifier = Modifier,
) {
    val isLoading = uiState is AuthUiState.Loading

    Column(
        modifier =
            modifier
                .fillMaxSize()
                .padding(horizontal = 32.dp),
        verticalArrangement = Arrangement.Center,
        horizontalAlignment = Alignment.CenterHorizontally,
    ) {
        Text(
            text = "Sign In",
            style = MaterialTheme.typography.headlineMedium,
            color = MaterialTheme.colorScheme.primary,
        )
        Spacer(modifier = Modifier.height(32.dp))
        OutlinedButton(
            onClick = onSignInWithGoogle,
            enabled = !isLoading,
            modifier = Modifier.fillMaxWidth(),
        ) {
            Icon(
                imageVector = GoogleIcon,
                contentDescription = null,
                tint = Color.Unspecified,
                modifier = Modifier.size(18.dp),
            )
            Spacer(modifier = Modifier.size(8.dp))
            Text(text = "Continue with Google")
        }
        Spacer(modifier = Modifier.height(20.dp))
        Row(verticalAlignment = Alignment.CenterVertically, modifier = Modifier.fillMaxWidth()) {
            HorizontalDivider(modifier = Modifier.weight(1f))
            Text(text = "  or  ", style = MaterialTheme.typography.bodySmall, color = MaterialTheme.colorScheme.outline)
            HorizontalDivider(modifier = Modifier.weight(1f))
        }
        Spacer(modifier = Modifier.height(20.dp))
        OutlinedTextField(
            value = loginForm.email,
            onValueChange = onEmailChange,
            label = { Text("Email") },
            singleLine = true,
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email),
            modifier = Modifier.fillMaxWidth(),
        )
        Spacer(modifier = Modifier.height(16.dp))
        OutlinedTextField(
            value = loginForm.password,
            onValueChange = onPasswordChange,
            label = { Text("Password") },
            singleLine = true,
            visualTransformation = PasswordVisualTransformation(),
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
            modifier = Modifier.fillMaxWidth(),
        )
        if (uiState is AuthUiState.Error) {
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                text = uiState.message,
                color = MaterialTheme.colorScheme.error,
                style = MaterialTheme.typography.bodySmall,
            )
        }
        Spacer(modifier = Modifier.height(24.dp))
        if (isLoading) {
            CircularProgressIndicator()
        } else {
            Button(
                onClick = onSignIn,
                enabled = loginForm.email.isNotBlank() && loginForm.password.isNotBlank(),
                modifier =
                    Modifier
                        .fillMaxWidth()
                        .testTag(LoginTestTags.SIGN_IN_BUTTON),
            ) {
                Text(text = "Sign In")
            }
        }
        Spacer(modifier = Modifier.height(16.dp))
        Row(horizontalArrangement = Arrangement.Center, modifier = Modifier.fillMaxWidth()) {
            Text(
                text = "Don't have an account?",
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onSurface,
            )
            TextButton(onClick = onNavigateToRegister) { Text(text = "Register") }
        }
    }
}

@Preview(showBackground = true)
@Composable
fun LoginScreenPreview() {
    TemplateTheme {
        LoginScreen(
            loginForm = LoginFormState(),
            uiState = AuthUiState.Idle,
            onEmailChange = {},
            onPasswordChange = {},
            onSignIn = {},
            onSignInWithGoogle = {},
            onNavigateToRegister = {},
        )
    }
}

@Preview(showBackground = true)
@Composable
fun LoginScreenErrorPreview() {
    TemplateTheme {
        LoginScreen(
            loginForm = LoginFormState(email = "a@b.com", password = "wrong"),
            uiState = AuthUiState.Error("Invalid credentials"),
            onEmailChange = {},
            onPasswordChange = {},
            onSignIn = {},
            onSignInWithGoogle = {},
            onNavigateToRegister = {},
        )
    }
}
