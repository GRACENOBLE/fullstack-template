package com.company.template.auth

import androidx.compose.ui.test.assertIsDisplayed
import androidx.compose.ui.test.junit4.createComposeRule
import androidx.compose.ui.test.onNodeWithText
import androidx.compose.ui.test.performClick
import androidx.test.ext.junit.runners.AndroidJUnit4
import com.company.template.ui.theme.TemplateTheme
import org.junit.Assert.assertTrue
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith

@RunWith(AndroidJUnit4::class)
class RegisterScreenTest {

    @get:Rule
    val composeTestRule = createComposeRule()

    @Test
    fun registerScreen_displaysCreateAccountHeading() {
        composeTestRule.setContent {
            TemplateTheme {
                RegisterScreen(
                    uiState = AuthUiState.Idle,
                    onRegister = { _, _, _, _ -> },
                    onNavigateToLogin = {},
                    onClearError = {}
                )
            }
        }
        composeTestRule.onNodeWithText("Create Account").assertIsDisplayed()
    }

    @Test
    fun registerScreen_displaysNameField() {
        composeTestRule.setContent {
            TemplateTheme {
                RegisterScreen(
                    uiState = AuthUiState.Idle,
                    onRegister = { _, _, _, _ -> },
                    onNavigateToLogin = {},
                    onClearError = {}
                )
            }
        }
        composeTestRule.onNodeWithText("Full Name").assertIsDisplayed()
    }

    @Test
    fun registerScreen_displaysErrorMessage() {
        composeTestRule.setContent {
            TemplateTheme {
                RegisterScreen(
                    uiState = AuthUiState.Error("Email already in use"),
                    onRegister = { _, _, _, _ -> },
                    onNavigateToLogin = {},
                    onClearError = {}
                )
            }
        }
        composeTestRule.onNodeWithText("Email already in use").assertIsDisplayed()
    }

    @Test
    fun registerScreen_displaysSignInLink() {
        composeTestRule.setContent {
            TemplateTheme {
                RegisterScreen(
                    uiState = AuthUiState.Idle,
                    onRegister = { _, _, _, _ -> },
                    onNavigateToLogin = {},
                    onClearError = {}
                )
            }
        }
        composeTestRule.onNodeWithText("Sign In").assertIsDisplayed()
    }

    @Test
    fun registerScreen_clickSignIn_invokesCallback() {
        var navigateCalled = false
        composeTestRule.setContent {
            TemplateTheme {
                RegisterScreen(
                    uiState = AuthUiState.Idle,
                    onRegister = { _, _, _, _ -> },
                    onNavigateToLogin = { navigateCalled = true },
                    onClearError = {}
                )
            }
        }
        composeTestRule.onNodeWithText("Sign In").performClick()
        assertTrue(navigateCalled)
    }
}
