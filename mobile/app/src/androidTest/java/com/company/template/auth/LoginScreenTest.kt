package com.company.template.auth

import androidx.compose.ui.test.assertIsDisplayed
import androidx.compose.ui.test.junit4.createComposeRule
import androidx.compose.ui.test.onNodeWithTag
import androidx.compose.ui.test.onNodeWithText
import androidx.compose.ui.test.performClick
import androidx.compose.ui.test.performTextInput
import androidx.test.ext.junit.runners.AndroidJUnit4
import com.company.template.ui.theme.TemplateTheme
import org.junit.Assert.assertTrue
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith

@RunWith(AndroidJUnit4::class)
class LoginScreenTest {
    @get:Rule
    val composeTestRule = createComposeRule()

    private var capturedEmail = ""
    private var capturedPassword = ""

    private fun setContent(
        loginForm: LoginFormState = LoginFormState(),
        uiState: AuthUiState = AuthUiState.Idle,
        onEmailChange: (String) -> Unit = { capturedEmail = it },
        onPasswordChange: (String) -> Unit = { capturedPassword = it },
        onSignIn: () -> Unit = {},
        onSignInWithGoogle: () -> Unit = {},
        onNavigateToRegister: () -> Unit = {},
    ) {
        composeTestRule.setContent {
            TemplateTheme {
                LoginScreen(
                    loginForm = loginForm,
                    uiState = uiState,
                    onEmailChange = onEmailChange,
                    onPasswordChange = onPasswordChange,
                    onSignIn = onSignIn,
                    onSignInWithGoogle = onSignInWithGoogle,
                    onNavigateToRegister = onNavigateToRegister,
                )
            }
        }
    }

    @Test
    fun loginScreen_displaysSignInHeading() {
        setContent()
        // The heading "Sign In" is displayed (there may also be a button with the same label)
        composeTestRule.onNodeWithText("Sign In", useUnmergedTree = false).assertIsDisplayed()
    }

    @Test
    fun loginScreen_displaysGoogleSignInButton() {
        setContent()
        composeTestRule.onNodeWithText("Continue with Google").assertIsDisplayed()
    }

    @Test
    fun loginScreen_googleButton_invokesCallback() {
        var called = false
        setContent(onSignInWithGoogle = { called = true })
        composeTestRule.onNodeWithText("Continue with Google").performClick()
        assertTrue(called)
    }

    @Test
    fun loginScreen_displaysErrorMessage() {
        setContent(uiState = AuthUiState.Error("Invalid credentials"))
        composeTestRule.onNodeWithText("Invalid credentials").assertIsDisplayed()
    }

    @Test
    fun loginScreen_displaysRegisterLink() {
        setContent()
        composeTestRule.onNodeWithText("Register").assertIsDisplayed()
    }

    @Test
    fun loginScreen_clickRegister_invokesCallback() {
        var navigateCalled = false
        setContent(onNavigateToRegister = { navigateCalled = true })
        composeTestRule.onNodeWithText("Register").performClick()
        assertTrue(navigateCalled)
    }

    @Test
    fun loginScreen_clickSignIn_invokesCallback() {
        var signInCalled = false
        setContent(
            loginForm = LoginFormState(email = "test@example.com", password = "secret"),
            onSignIn = { signInCalled = true },
        )
        // Use the test tag to target the button specifically, not the heading
        composeTestRule.onNodeWithTag(LoginTestTags.SIGN_IN_BUTTON).performClick()
        assertTrue(signInCalled)
    }

    @Test
    fun loginScreen_typingEmail_updatesCallback() {
        var updated = ""
        setContent(onEmailChange = { updated = it })
        composeTestRule.onNodeWithText("Email").performTextInput("test@example.com")
        assertTrue(updated.isNotEmpty())
    }

    @Test
    fun loginScreen_typingPassword_updatesCallback() {
        var updated = ""
        setContent(onPasswordChange = { updated = it })
        composeTestRule.onNodeWithText("Password").performTextInput("secret123")
        assertTrue(updated.isNotEmpty())
    }
}
