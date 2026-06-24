package com.company.template.onboarding

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
class OnboardingScreenTest {

    @get:Rule
    val composeTestRule = createComposeRule()

    @Test
    fun onboardingScreen_displaysWelcomeHeading() {
        composeTestRule.setContent {
            TemplateTheme {
                OnboardingScreen(onGetStarted = {})
            }
        }
        composeTestRule.onNodeWithText("Welcome").assertIsDisplayed()
    }

    @Test
    fun onboardingScreen_displaysGetStartedButton() {
        composeTestRule.setContent {
            TemplateTheme {
                OnboardingScreen(onGetStarted = {})
            }
        }
        composeTestRule.onNodeWithText("Get Started").assertIsDisplayed()
    }

    @Test
    fun onboardingScreen_clickGetStarted_invokesCallback() {
        var callbackInvoked = false
        composeTestRule.setContent {
            TemplateTheme {
                OnboardingScreen(onGetStarted = { callbackInvoked = true })
            }
        }
        composeTestRule.onNodeWithText("Get Started").performClick()
        assertTrue(callbackInvoked)
    }
}
