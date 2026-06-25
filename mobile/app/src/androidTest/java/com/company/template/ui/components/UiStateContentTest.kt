package com.company.template.ui.components

import androidx.compose.material3.Text
import androidx.compose.ui.test.assertIsDisplayed
import androidx.compose.ui.test.junit4.createComposeRule
import androidx.compose.ui.test.onNodeWithTag
import androidx.compose.ui.test.onNodeWithText
import androidx.compose.ui.test.performClick
import androidx.test.ext.junit.runners.AndroidJUnit4
import com.company.template.ui.state.UiState
import com.company.template.ui.theme.TemplateTheme
import org.junit.Assert.assertTrue
import org.junit.Rule
import org.junit.Test
import org.junit.runner.RunWith

@RunWith(AndroidJUnit4::class)
class UiStateContentTest {
    @get:Rule
    val composeTestRule = createComposeRule()

    @Test
    fun loadingState_showsCircularProgressIndicator() {
        composeTestRule.setContent {
            TemplateTheme {
                UiStateContent(
                    state = UiState.Loading,
                    content = { _: Unit -> },
                )
            }
        }

        composeTestRule.onNodeWithTag("loading_indicator").assertIsDisplayed()
    }

    @Test
    fun errorState_showsErrorMessage() {
        composeTestRule.setContent {
            TemplateTheme {
                UiStateContent(
                    state = UiState.Error("Something went wrong"),
                    content = { _: Unit -> },
                )
            }
        }

        composeTestRule.onNodeWithText("Something went wrong").assertIsDisplayed()
    }

    @Test
    fun errorState_withRetry_showsRetryButton() {
        var retryClicked = false

        composeTestRule.setContent {
            TemplateTheme {
                UiStateContent(
                    state = UiState.Error("Network error"),
                    onRetry = { retryClicked = true },
                    content = { _: Unit -> },
                )
            }
        }

        composeTestRule.onNodeWithText("Retry").assertIsDisplayed()
        composeTestRule.onNodeWithText("Retry").performClick()
        assertTrue(retryClicked)
    }

    @Test
    fun errorState_withoutRetry_doesNotShowRetryButton() {
        composeTestRule.setContent {
            TemplateTheme {
                UiStateContent(
                    state = UiState.Error("Network error"),
                    onRetry = null,
                    content = { _: Unit -> },
                )
            }
        }

        composeTestRule.onNodeWithText("Retry").assertDoesNotExist()
    }

    @Test
    fun successState_rendersContentLambda() {
        composeTestRule.setContent {
            TemplateTheme {
                UiStateContent(
                    state = UiState.Success("Hello from success"),
                ) { data ->
                    Text(text = data)
                }
            }
        }

        composeTestRule.onNodeWithText("Hello from success").assertIsDisplayed()
    }

    @Test
    fun idleState_rendersNothing() {
        composeTestRule.setContent {
            TemplateTheme {
                UiStateContent(
                    state = UiState.Idle,
                    content = { _: Unit -> Text(text = "should not appear") },
                )
            }
        }

        composeTestRule.onNodeWithText("should not appear").assertDoesNotExist()
        composeTestRule.onNodeWithTag("loading_indicator").assertDoesNotExist()
    }
}
