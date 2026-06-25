package com.company.template

import androidx.compose.ui.test.assertIsDisplayed
import androidx.compose.ui.test.junit4.createComposeRule
import androidx.compose.ui.test.onNodeWithText
import androidx.test.ext.junit.runners.AndroidJUnit4
import com.company.template.ui.theme.TemplateTheme
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith

@RunWith(AndroidJUnit4::class)
class GreetingTest {
    @get:Rule
    val composeTestRule = createComposeRule()

    @Test
    fun greeting_displaysName() {
        composeTestRule.setContent {
            TemplateTheme {
                Greeting(name = "World")
            }
        }
        composeTestRule.onNodeWithText("Hello World!").assertIsDisplayed()
    }

    @Test
    fun greeting_displaysAndroid() {
        composeTestRule.setContent {
            TemplateTheme {
                Greeting(name = "Android")
            }
        }
        composeTestRule.onNodeWithText("Hello Android!").assertIsDisplayed()
    }
}
