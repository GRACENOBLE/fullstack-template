package com.company.template.settings

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
class SettingsScreenTest {
    @get:Rule
    val composeTestRule = createComposeRule()

    @Test
    fun settingsScreen_displaysDisplayNameAndEmail() {
        composeTestRule.setContent {
            TemplateTheme {
                SettingsScreen(
                    displayName = "Alice Example",
                    email = "alice@example.com",
                    onSignOut = {},
                )
            }
        }
        composeTestRule.onNodeWithText("Alice Example").assertIsDisplayed()
        composeTestRule.onNodeWithText("alice@example.com").assertIsDisplayed()
    }

    @Test
    fun settingsScreen_clickSignOut_invokesCallback() {
        var signOutCalled = false
        composeTestRule.setContent {
            TemplateTheme {
                SettingsScreen(
                    displayName = "Alice Example",
                    email = "alice@example.com",
                    onSignOut = { signOutCalled = true },
                )
            }
        }
        composeTestRule.onNodeWithText("Sign out").performClick()
        assertTrue(signOutCalled)
    }

    @Test
    fun settingsScreen_nullDisplayNameAndEmail_showsDashes() {
        composeTestRule.setContent {
            TemplateTheme {
                SettingsScreen(
                    displayName = null,
                    email = null,
                    onSignOut = {},
                )
            }
        }
        composeTestRule.onNodeWithText("—").assertIsDisplayed()
    }
}
