package com.company.template.auth

import androidx.compose.ui.test.assertIsDisplayed
import androidx.compose.ui.test.junit4.createComposeRule
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

    private fun setContent(
        uiState: AuthUiState = AuthUiState.Idle,
        onSignIn: (String, String) -> Unit = { _, _ -> },
        onSignInWithGoogle: () -> Unit = {},
        onNavigateToRegister: () -> Unit = {},
        onClearError: () -> Unit = {}
    ) {
        composeTestRule.setContent {
            TemplateTheme {
                LoginScreen(
                    uiState = uiState,
                    onSignIn = onSignIn,
                    onSignInWithGoogle = onSignInWithGoogle,
                    onNavigateToRegister = onNavigateToRegister,
                    onClearError = onClearError
                )
            }
        }
    }

    @Test
    fun loginScreen_displaysSignInHeading() {
        setContent()
        composeTestRule.onNodeWithText("Sign In").assertIsDisplayed()
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
    fun loginScreen_typingEmailAndPassword_thenSignInInvoked() {
        var signInEmail = ""
        var signInPassword = ""
        setContent(onSignIn = { e, p -> signInEmail = e; signInPassword = p })
        composeTestRule.onNodeWithText("Email").performTextInput("test@example.com")
        composeTestRule.onNodeWithText("Password").performTextInput("secret123")
        composeTestRule.onNodeWithText("Sign In").performClick()
        assertTrue(signInEmail == "test@example.com")
        assertTrue(signInPassword == "secret123")
    }
}
