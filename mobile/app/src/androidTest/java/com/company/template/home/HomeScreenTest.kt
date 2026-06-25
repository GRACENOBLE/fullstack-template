package com.company.template.home

import androidx.compose.ui.test.assertDoesNotExist
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
class HomeScreenTest {
    @get:Rule
    val composeTestRule = createComposeRule()

    @Test
    fun homeScreen_displaysWelcomeHeading() {
        composeTestRule.setContent {
            TemplateTheme {
                HomeScreen(displayName = "Alice", onSignOut = {})
            }
        }
        composeTestRule.onNodeWithText("Welcome back!").assertIsDisplayed()
    }

    @Test
    fun homeScreen_displaysDisplayName() {
        composeTestRule.setContent {
            TemplateTheme {
                HomeScreen(displayName = "Alice", onSignOut = {})
            }
        }
        composeTestRule.onNodeWithText("Alice").assertIsDisplayed()
    }

    @Test
    fun homeScreen_displaysSignOutButton() {
        composeTestRule.setContent {
            TemplateTheme {
                HomeScreen(displayName = "Alice", onSignOut = {})
            }
        }
        composeTestRule.onNodeWithText("Sign Out").assertIsDisplayed()
    }

    @Test
    fun homeScreen_clickSignOut_invokesCallback() {
        var signOutCalled = false
        composeTestRule.setContent {
            TemplateTheme {
                HomeScreen(displayName = "Alice", onSignOut = { signOutCalled = true })
            }
        }
        composeTestRule.onNodeWithText("Sign Out").performClick()
        assertTrue(signOutCalled)
    }

    @Test
    fun homeScreen_emptyDisplayName_doesNotShowNameText() {
        composeTestRule.setContent {
            TemplateTheme {
                HomeScreen(displayName = "", onSignOut = {})
            }
        }
        composeTestRule.onNodeWithText("Welcome back!").assertIsDisplayed()
        // An empty string should not produce a visible name node
        composeTestRule.onNodeWithText("").assertDoesNotExist()
    }
}
