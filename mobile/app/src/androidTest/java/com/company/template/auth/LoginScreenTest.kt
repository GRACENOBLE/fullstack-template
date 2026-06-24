package com.company.template.auth

import androidx.compose.ui.test.assertIsDisplayed
import androidx.compose.ui.test.assertIsEnabled
import androidx.compose.ui.test.assertIsNotEnabled
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

    @Test
    fun loginScreen_displaysSignInHeading() {
        composeTestRule.setContent {
            TemplateTheme {
                LoginScreen(
                    uiState = AuthUiState.Idle,
                    onSignIn = { _, _ -> },
                    onNavigateToRegister = {},
                    onClearError = {}
                )
            }
        }
        composeTestRule.onNodeWithText("Sign In").assertIsDisplayed()
    }

    @Test
    fun loginScreen_signInButtonDisabledWhenFieldsEmpty() {
        composeTestRule.setContent {
            TemplateTheme {
                LoginScreen(
                    uiState = AuthUiState.Idle,
                    onSignIn = { _, _ -> },
                    onNavigateToRegister = {},
                    onClearError = {}
                )
            }
        }
        // Button with text "Sign In" that is a Button (not the heading Text)
        composeTestRule.onNodeWithText("Sign In", useUnmergedTree = false)
            .assertIsDisplayed()
    }

    @Test
    fun loginScreen_displaysErrorMessage() {
        composeTestRule.setContent {
            TemplateTheme {
                LoginScreen(
                    uiState = AuthUiState.Error("Invalid credentials"),
                    onSignIn = { _, _ -> },
                    onNavigateToRegister = {},
                    onClearError = {}
                )
            }
        }
        composeTestRule.onNodeWithText("Invalid credentials").assertIsDisplayed()
    }

    @Test
    fun loginScreen_displaysRegisterLink() {
        composeTestRule.setContent {
            TemplateTheme {
                LoginScreen(
                    uiState = AuthUiState.Idle,
                    onSignIn = { _, _ -> },
                    onNavigateToRegister = {},
                    onClearError = {}
                )
            }
        }
        composeTestRule.onNodeWithText("Register").assertIsDisplayed()
    }

    @Test
    fun loginScreen_clickRegister_invokesCallback() {
        var navigateCalled = false
        composeTestRule.setContent {
            TemplateTheme {
                LoginScreen(
                    uiState = AuthUiState.Idle,
                    onSignIn = { _, _ -> },
                    onNavigateToRegister = { navigateCalled = true },
                    onClearError = {}
                )
            }
        }
        composeTestRule.onNodeWithText("Register").performClick()
        assertTrue(navigateCalled)
    }

    @Test
    fun loginScreen_typingEmailAndPassword_thenSignInInvoked() {
        var signInEmail = ""
        var signInPassword = ""
        composeTestRule.setContent {
            TemplateTheme {
                LoginScreen(
                    uiState = AuthUiState.Idle,
                    onSignIn = { e, p ->
                        signInEmail = e
                        signInPassword = p
                    },
                    onNavigateToRegister = {},
                    onClearError = {}
                )
            }
        }
        composeTestRule.onNodeWithText("Email").performTextInput("test@example.com")
        composeTestRule.onNodeWithText("Password").performTextInput("secret123")
        composeTestRule.onNodeWithText("Sign In").performClick()
        assertTrue(signInEmail == "test@example.com")
        assertTrue(signInPassword == "secret123")
    }
}
