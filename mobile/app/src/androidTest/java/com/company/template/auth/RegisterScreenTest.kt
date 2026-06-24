package com.company.template.auth

import androidx.compose.ui.test.assertIsDisplayed
import androidx.compose.ui.test.junit4.createComposeRule
import androidx.compose.ui.test.onNodeWithTag
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

    private fun setContent(
        registerForm: RegisterFormState = RegisterFormState(),
        uiState: AuthUiState = AuthUiState.Idle,
        onNameChange: (String) -> Unit = {},
        onEmailChange: (String) -> Unit = {},
        onPasswordChange: (String) -> Unit = {},
        onConfirmPasswordChange: (String) -> Unit = {},
        onRegister: () -> Unit = {},
        onNavigateToLogin: () -> Unit = {},
    ) {
        composeTestRule.setContent {
            TemplateTheme {
                RegisterScreen(
                    registerForm = registerForm,
                    uiState = uiState,
                    onNameChange = onNameChange,
                    onEmailChange = onEmailChange,
                    onPasswordChange = onPasswordChange,
                    onConfirmPasswordChange = onConfirmPasswordChange,
                    onRegister = onRegister,
                    onNavigateToLogin = onNavigateToLogin,
                )
            }
        }
    }

    @Test
    fun registerScreen_displaysCreateAccountHeading() {
        setContent()
        // Heading text (there may also be a button with the same label; use useUnmergedTree if needed)
        composeTestRule.onNodeWithText("Create Account", useUnmergedTree = false).assertIsDisplayed()
    }

    @Test
    fun registerScreen_displaysNameField() {
        setContent()
        composeTestRule.onNodeWithText("Full Name").assertIsDisplayed()
    }

    @Test
    fun registerScreen_displaysErrorMessage() {
        setContent(uiState = AuthUiState.Error("Email already in use"))
        composeTestRule.onNodeWithText("Email already in use").assertIsDisplayed()
    }

    @Test
    fun registerScreen_displaysSignInLink() {
        setContent()
        composeTestRule.onNodeWithText("Sign In").assertIsDisplayed()
    }

    @Test
    fun registerScreen_clickSignIn_invokesCallback() {
        var navigateCalled = false
        setContent(onNavigateToLogin = { navigateCalled = true })
        composeTestRule.onNodeWithText("Sign In").performClick()
        assertTrue(navigateCalled)
    }

    @Test
    fun registerScreen_clickCreateAccount_invokesCallback() {
        var registerCalled = false
        setContent(
            registerForm = RegisterFormState(
                name = "Alice",
                email = "alice@example.com",
                password = "pass",
                confirmPassword = "pass",
            ),
            onRegister = { registerCalled = true },
        )
        // Use the test tag to target the button specifically, not the heading
        composeTestRule.onNodeWithTag(RegisterTestTags.CREATE_ACCOUNT_BUTTON).performClick()
        assertTrue(registerCalled)
    }
}
